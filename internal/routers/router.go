package routers

import (
	"quick-im-demo/internal/controllers"
	"quick-im-demo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	registerRouter := r.Group("register")
	{
		//中间件直接当函数写，无需套用结构体和*c.Context
		registerRouter.Use(middleware.DefaultHandler())
		//注意注册时结构体参数不需要指针类型
		registerRouter.GET("/", controllers.RegisterController{}.Register)
	}

	r.Run(":8080")
}
