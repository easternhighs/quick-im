package controllers

import "github.com/gin-gonic/gin"

type DefaultController struct {
}

func (dc DefaultController) Index(c *gin.Context) {
	c.String(200, "Welcome to the default controller!")
}
