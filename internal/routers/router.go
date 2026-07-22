package routers

import (
	"net/http"
	"quick-im-demo/internal/controllers"
	"quick-im-demo/internal/middlewares"
	"quick-im-demo/internal/ws"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.LoadHTMLGlob("internal/templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/auth/login") })

	authRouter := r.Group("/auth")
	{
		authRouter.Use(middlewares.DefaultHandler())
		authRouter.GET("/register", controllers.RegisterController{}.DoRegister)
		authRouter.POST("/register", controllers.RegisterController{}.WriteRegisterInfo)
		authRouter.GET("/login", controllers.LoginController{}.ShowLogin)
		authRouter.POST("/login", controllers.LoginController{}.DoLogin)
		authRouter.POST("/logout", controllers.LoginController{}.Logout)
	}

	hub := ws.NewHub()
	r.GET("/ws", (&controllers.WebSocketController{Hub: hub}).Connect)

	imRouter := r.Group("/im")
	{
		imRouter.Use(middlewares.AuthHandler())
		imRouter.GET("", controllers.IMController{}.Dashboard)
	}

	_ = r.Run(":8000")
}
