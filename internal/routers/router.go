package routers

import (
	"context"
	"net/http"
	"quick-im-demo/internal/cache"
	"quick-im-demo/internal/config"
	"quick-im-demo/internal/controllers"
	"quick-im-demo/internal/middlewares"
	"quick-im-demo/internal/models/user"
	"quick-im-demo/internal/realtime"
	"quick-im-demo/internal/ws"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	conf := config.ReadConfig()
	r := gin.Default()
	r.LoadHTMLGlob("internal/templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/auth/login") })

	authRouter := r.Group("/auth")
	authRouter.Use(middlewares.DefaultHandler())
	authRouter.GET("/register", controllers.RegisterController{}.DoRegister)
	authRouter.POST("/register", controllers.RegisterController{}.WriteRegisterInfo)
	authRouter.GET("/login", controllers.LoginController{}.ShowLogin)
	authRouter.POST("/login", controllers.LoginController{}.DoLogin)
	authRouter.POST("/logout", controllers.LoginController{}.Logout)

	hub, rdb := ws.NewHub(), cache.NewRedis()
	channel := conf.Redis.Channel
	go realtime.SubscribeMessages(context.Background(), rdb, channel, hub)
	r.GET("/ws", (&controllers.WebSocketController{Hub: hub}).Connect)

	imRouter := r.Group("/im")
	imRouter.Use(middlewares.AuthHandler())
	imRouter.GET("", controllers.IMController{}.Dashboard)

	messageController := &controllers.MessageController{DB: user.DB, Hub: hub, Redis: rdb, Channel: channel}
	apiRouter := r.Group("/api")
	apiRouter.Use(middlewares.AuthHandler())
	apiRouter.POST("/messages", messageController.Send)
	apiRouter.GET("/messages/history", messageController.History)
	apiRouter.POST("/messages/read", messageController.MarkRead)
	apiRouter.GET("/unread", messageController.Unread)
	apiRouter.GET("/conversations", messageController.Conversations)

	_ = r.Run(conf.App.HTTPAddr)
}
