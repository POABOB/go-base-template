package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/POABOB/go-base-template/handler"
	hystrix_micro "github.com/POABOB/go-base-template/plugin/hystrix"
	base "github.com/POABOB/go-base-template/protos/base"
	common "github.com/POABOB/utils"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/opentracing/opentracing-go"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentrace_micro "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 1. 註冊中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		// 註冊中心地址
		options.Addrs = []string{
			//如果放到docker-compose 也可以是服務名稱
			// "node4",
			"localhost:8500",
		}
	})

	// 2. 配置中心，存放經常變動的變量
	// 使用之前請先在http://localhost:8500的Key/Value頁面新增/micro/config/mysql配置
	// /micro/config/mysql
	// {
	//   "host": "127.0.0.1",
	//   "user": "root",
	//   "password": "root",
	//   "database": "PaaS",
	//   "port": 3306
	// }
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/config")
	if err != nil {
		common.Error(err)
	}

	// 3. 使用配置中心連接 mysql
	// 如果遇到 Error 1049: Unknown database 'PaaS' 錯誤
	// 請先連接 mysql 並且建立好資料庫
	// SQL指令： CREATE DATABASE PaaS CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	mysqlConfig := common.GetMysqlFromConsul(consulConfig, "mysql")
	// 初始化 DB
	db, err := gorm.Open("mysql", mysqlConfig.User+":"+mysqlConfig.Password+"@/"+mysqlConfig.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Fatal(err)
	}
	common.Info("連接 mysql 成功")
	// 如果 db 不使用，就把它關掉，defer 最後執行
	defer db.Close()
	// 禁止重複建立資料表
	db.SingularTable(true)

	// 4. 新增鏈路追蹤
	t, io, err := common.NewTracer("base", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	// 如果 io 不使用，就把它關掉，defer 最後執行
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 5. 新增熔斷器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	// 啟動監聽程式

	// 6. 新增日誌中心
	// 1）需要程式日誌打入到日誌文件中
	// 2）在程式中新增 filebeat.yml 文件
	// 3) 啟動 filebeat，啟動命令 ./filebeat -e -c filebeat.yml
	common.Info("新增 日誌系统 ！")

	go func() {
		// http://192.168.1.105:9092/turbine/turbine.stream
		// 看板訪問地址 http://127.0.0.1:9002/hystrix，url後面一定要有 /hystrix
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9092"), hystrixStreamHandler)
		common.Info("熔斷")
		fmt.Println("333")
		if err != nil {
			common.Error(err)
		}
	}()

	// // 7. 新增監控
	common.PrometheusBoot(9192)

	// 新增服務
	service := micro.NewService(
		// 名稱、版本
		micro.Name("base"),
		micro.Version("latest"),

		// 新增註冊中心
		micro.Registry(consul),

		// 新增鏈路追蹤
		// 紀錄被訪問的請求
		micro.WrapHandler(opentrace_micro.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 紀錄我們訪問別人的請求
		micro.WrapClient(opentrace_micro.NewClientWrapper(opentracing.GlobalTracer())),
		// 只作為客戶端時候起作用
		micro.WrapClient(hystrix_micro.NewClientHystrixWrapper()),

		// 新增限流，1000 req/s 超出就不處理
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
	)

	// 初始化
	service.Init()

	// 註冊句柄，可以快速操作已開發服務
	base.RegisterBaseHandler(service.Server(), new(handler.Base))

	// 啟動服務
	if err := service.Run(); err != nil {
		// 輸出失敗訊息
		common.Fatal(err)
	}
}
