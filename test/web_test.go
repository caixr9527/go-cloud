package test

import (
	"fmt"
	"github.com/caixr9527/go-cloud"
	"github.com/caixr9527/go-cloud/common"
	"github.com/caixr9527/go-cloud/web"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	engine := cloud.Default()
	engine.Use(func(context *web.Context) {
		fmt.Println("Global before")
		context.Next()
		fmt.Println("Global after")
	})
	handle := engine.Handle()
	group := handle.Group("user")
	group.Use(func(context *web.Context) {
		fmt.Println("group before")
		context.Next()
		fmt.Println("group after")
	})
	group.GET("/hello/:id/:name", func(context *web.Context) {
		fmt.Println("user middle")
	}, func(context *web.Context) {
		fmt.Println("user middle1")
	}, func(context *web.Context) {
		fmt.Println("good")
		data := make([]any, 0)
		data = append(data, context.PathVariable("id"), context.PathVariable("name"), "hello,go_cloud")
		r := &common.R{
			Success: true,
			Code:    http.StatusOK,
			Data:    data,
			Msg:     "",
		}
		context.JSON(http.StatusOK, r)
	})
	handler := handle.Group("/order")
	handler.GET("/hello", func(context *web.Context) {
		fmt.Println("order middle")
	}, func(context *web.Context) {
		context.JSON(http.StatusOK, "hello,go_cloud2")
	})
	handler.POST("/hello", func(context *web.Context) {
		fmt.Println("post order middle")
		context.Data = "GGGGG"
	}, func(context *web.Context) {
		context.JSON(http.StatusOK, "post hello,go_cloud2", context.Data)
	})
	engine.Run(":8111")
}
