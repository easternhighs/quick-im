package middlewares

import (
	"net/http"
	"quick-im-demo/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Authorization: Bearer 是一种在 HTTP 请求头部中用于传递访问令牌（Access Token）的常见格式，
		//用于在客户端和服务器之间进行身份验证和授权操作。
		//Authorization: Bearer your_access_token
		tokenText := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenText == "" {
			tokenText, _ = c.Cookie("quick_im_token")
		}

		token, err := utils.VerifyToken(tokenText)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		username, ok := claims["username"].(string)
		if !ok || username == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
