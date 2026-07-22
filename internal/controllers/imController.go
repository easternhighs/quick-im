package controllers

import "github.com/gin-gonic/gin"

type IMController struct{}

func (IMController) Dashboard(c *gin.Context) {
	c.HTML(200, "im.html", gin.H{"username": c.GetString("username")})
}
