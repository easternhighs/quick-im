package middlewares

import "github.com/gin-gonic/gin"

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加认证逻辑，例如检查请求头中的 token 是否有效
	}
}
