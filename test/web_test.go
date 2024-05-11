package test

import (
	"fmt"
	go_cloud "github.com/caixr9527/go-cloud"
	"github.com/caixr9527/go-cloud/web"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	engine := go_cloud.New()
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
		context.JSON(http.StatusOK, "hello,go_cloud", context.PathVariable("id"), context.PathVariable("name"))
	})
	handler := handle.Group("/order")
	handler.GET("/hello", func(context *web.Context) {
		fmt.Println("order middle")
	}, func(context *web.Context) {
		context.JSON(http.StatusOK, "hello,go_cloud2")
	})
	engine.Run(":8111")
}
