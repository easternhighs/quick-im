package middleware

import "github.com/gin-gonic/gin"

func DefaultHandler() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
