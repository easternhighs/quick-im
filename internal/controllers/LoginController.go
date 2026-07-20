package controllers

import "github.com/gin-gonic/gin"

type LoginController struct {
	// DefaultController
}

func (lc LoginController) DoLogin(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}
