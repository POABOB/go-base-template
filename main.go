package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/PPOABOB/go-base-template/handler"
	hystrix2 "github.com/PPOABOB/go-base-template/plugin/hystrix"
	base "github.com/PPOABOB/go-base-template/protos/base"

	"github.com/POABOB/common"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	log "github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// 1. 註冊中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			//如果放到docker-compose 也可以是服務名稱
			"localhost:8500",
		}
	})
	// 2. 配置中心，存放經常變動的變量
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/config")
	if err != nil {
		common.Error(err)
	}

	// 3. 使用配置中心連接 mysql
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	// 初始化 DB
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Fatal(err)
	}
	fmt.Println("連接mysql 成功")
	common.Info("連接mysql 成功")
	defer db.Close()
	// 禁止複表
	db.SingularTable(true)

	// 4. 新增鏈路追蹤
	t, io, err := common.NewTracer("base", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 5. 新增熔斷器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	// 6. 新增日誌中心
	// 1）需要程式日誌打入到日誌文件中
	// 2）在程式中新增 filebeat.yml 文件
	// 3) 啟動 filebeat，啟動命令 ./filebeat -e -c filebeat.yml
	common.Info("新增 日誌系统 ！")

	// 啟動監聽程式
	go func() {
		// http://192.168.0.112:9092/turbine/turbine.stream
		// 看板訪問地址 http://127.0.0.1:9002/hystrix，url後面一定要有 /hystrix
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9092"), hystrixStreamHandler)
		fmt.Println("333")
		if err != nil {
			fmt.Println(err)
		}
	}()

	// 7. 新增監控
	common.PrometheusBoot(9192)

	// 新增服務
	service := micro.NewService(
		micro.Name("base"),
		micro.Version("latest"),
		// 新增註冊中心
		micro.Registry(consul),
		// 新增鏈路追蹤
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 只作為客戶端時候起作用
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// 新增限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
	)

	// 初始化
	service.Init()

	// 註冊句柄，可以快速操作已開發服務
	base.RegisterBaseHandler(service.Server(), new(handler.Base))

	// 啟動服務
	if err := service.Run(); err != nil {
		// 輸出失敗訊息
		log.Fatal(err)
	}
}
