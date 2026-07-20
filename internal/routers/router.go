package routers

import (
	"net/http"
	"quick-im-demo/internal/controllers"
	"quick-im-demo/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	// LoadHTMLGlob uses the process working directory. The server is started
	// from the project root, so templates must be addressed from that root.
	r.LoadHTMLGlob("internal/templates/*")
	//将/static路径映射到本地的static目录下，方便访问静态资源
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth/register")
	})

	authRouter := r.Group("/auth")
	{
		//中间件直接当函数写，无需套用结构体和*c.Context
		authRouter.Use(middlewares.DefaultHandler())
		//注意注册时结构体参数不需要指针类型
		authRouter.GET("/register", controllers.RegisterController{}.DoRegister)
		authRouter.POST("/register", controllers.RegisterController{}.WriteRegisterInfo)

		//登录
		authRouter.POST("/login", controllers.LoginController{}.DoLogin)
	}

	_ = r.Run(":8000")
}
