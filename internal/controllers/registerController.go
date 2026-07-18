package controllers

import "github.com/gin-gonic/gin"

type RegisterController struct {
	// DefaultController
}

func (rc RegisterController) Register(c *gin.Context) {
	c.HTML(200, "admin.html", gin.H{})
}
