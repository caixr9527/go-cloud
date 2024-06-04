package test

import (
	"fmt"
	"github.com/caixr9527/go-cloud"
	"github.com/caixr9527/go-cloud/auth"
	"github.com/caixr9527/go-cloud/common"
	logger "github.com/caixr9527/go-cloud/log"
	"github.com/caixr9527/go-cloud/web"
	"log"
	"net/http"
	"testing"
	"time"
)

type User struct {
	Name      string   `xml:"name" json:"name" `
	Age       int      `xml:"age" json:"age" validate:"required,max=50,min=18"`
	Addresses []string `xml:"addresses" json:"addresses"`
}

func TestRun(t *testing.T) {

	options := &web.Options{}
	engine := cloud.New(options)
	//engine := cloud.Default()
	engine.Use(auth.Token)
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
		time.Sleep(time.Duration(10) * time.Second)
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
	group.GET("/hello2/:locationId/:username", func(context *web.Context) {
		fmt.Println("user middle")
	}, func(context *web.Context) {
		fmt.Println("user middle1")
	}, func(context *web.Context) {
		fmt.Println("good")
		data := make([]any, 0)
		data = append(data, context.PathVariable("locationId"), context.PathVariable("username"), "hello,go_cloud")
		r := &common.R{
			Success: true,
			Code:    http.StatusOK,
			Data:    data,
			Msg:     "",
		}
		context.JSON(http.StatusOK, r)
	})
	orderGroup := handle.Group("/order")
	orderGroup.GET("/hello", func(context *web.Context) {
		fmt.Println("order middle")
	}, func(context *web.Context) {
		context.JSON(http.StatusOK, "hello,go_cloud2")
	})
	orderGroup.POST("/hello", func(context *web.Context) {
		fmt.Println("post order middle")
		context.Data = "GGGGG"
	}, func(context *web.Context) {
		context.JSON(http.StatusOK, "post hello,go_cloud2", context.Data)
	})
	orderGroup.GET("/template", func(context *web.Context) {
		user := &User{Name: "caixiaorongtemplate"}
		context.ParseTemplate("login.html", user)
	})

	orderGroup.GET("/htmlTemplate", func(context *web.Context) {
		user := &User{Name: "caixiaoronghtmlTemplate"}
		context.ParseTemplates("index.html", user, "tpl/index.html")
	})

	orderGroup.GET("/html", func(context *web.Context) {
		key := context.Query("key")
		ids := context.QueryArray("ids")
		fmt.Println(ids)
		fmt.Println(context.QueryMap())
		keyDefault := context.QueryDefault("keyDefault", "keyDefault")
		context.ToHTML(200, "<h1>"+key+"</h1>\n"+"<h1>"+keyDefault+"</h1>\n")
	})

	orderGroup.GET("/fileDownload", func(context *web.Context) {
		context.FileDownload("tpl/1.xlsx")
	})

	orderGroup.POST("/fileDownload", func(context *web.Context) {
		context.FileDownloadWithFilename("tpl/1.xlsx", "aaa.xlsx")
	})

	orderGroup.GET("/fileFromFS", func(context *web.Context) {
		context.FileFromFS("1.xlsx", http.Dir("tpl"))
	})

	orderGroup.POST("/postForm", func(context *web.Context) {
		id, err := context.PostForm("id")
		if err != nil {
			log.Println(err)
		}
		name, err := context.PostForm("name")
		if err != nil {
			log.Println(err)
		}
		age, err := context.PostForm("age")
		if err != nil {
			log.Println(err)
		}
		context.JSON(http.StatusOK, id, name, age)
		//fmt.Println(context.PostFormArray("id"))
		//context.JSON(http.StatusOK, context.PostFormMap())
	})

	orderGroup.POST("/postForm2", func(context *web.Context) {
		id, err := context.PostForm("id")
		if err != nil {
			log.Println(err)
		}
		name, err := context.PostForm("name")
		if err != nil {
			log.Println(err)
		}
		age, err := context.PostForm("age")
		if err != nil {
			log.Println(err)
		}
		context.JSON(http.StatusOK, id, name, age)
		//fmt.Println(context.PostFormArray("id"))
		//context.JSON(http.StatusOK, context.PostFormMap())
	})

	orderGroup.POST("/postFormFile", func(context *web.Context) {
		//file := context.FormFile("file")
		files, _ := context.FormFiles("file")
		for _, file := range files {
			context.UploadedFile(file, "F:\\workspace\\personal\\upload\\"+file.Filename)
		}

		context.JSON(http.StatusOK)
	})

	orderGroup.POST("/bind", func(context *web.Context) {
		user := &User{}
		//var str string
		err := context.Bind(user)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, user)
	})

	orderGroup.GET("/bind2", func(context *web.Context) {
		user := &User{}
		//var str string
		err := context.BindQuery(user)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, user)
	})

	orderGroup.GET("/bind3", func(context *web.Context) {
		user := &User{}
		//var str string
		err := context.Bind(user)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, user)
	})

	orderGroup.GET("/bind4", func(context *web.Context) {
		//user := &User{}
		users := make([]User, 0)
		//var str string
		err := context.Bind(&users)
		if err != nil {
			logger.Log.Error(err.Error())
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		logger.Log.Debug("Debug")
		logger.Log.Info("Info")
		logger.Log.Warn("Warn")
		logger.Log.Error("Error")
		context.JSON(http.StatusOK, users)
	})
	config := &auth.JwtToken{
		Alg:          "HS256",
		TokenTimeout: 30 * time.Second,
		RefreshKey:   []byte("refreshKey"),
		Key:          []byte("go_cloud"),
	}
	claims := map[string]any{"userId": 1}
	token, _ := config.CreateToken(claims)
	fmt.Println(token.Token)
	fmt.Println(token.RefreshToken)
	engine.Run()
}
