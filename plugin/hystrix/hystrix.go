package hystrix

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/v3/client"
)

type clientWrapper struct {
	client.Client
}

// 熔斷邏輯
func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// 正常運作
		fmt.Println(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(e error) error {
		// 走熔斷邏輯，每個服務不一樣
		fmt.Println(req.Service() + "." + req.Endpoint() + "的熔斷邏輯")
		return e
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
