package controllers

import (
	"net/http"
	"quick-im-demo/internal/utils"
	"quick-im-demo/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

type WebSocketController struct{ Hub *ws.Hub }

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (wc *WebSocketController) Connect(c *gin.Context) {
	tokenText, _ := c.Cookie("quick_im_token")
	if tokenText == "" {
		tokenText = c.Query("token")
	}
	token, err := utils.VerifyToken(tokenText)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	wc.Hub.Add(username, conn)
	defer wc.Hub.Remove(username, conn)
	_ = conn.WriteJSON(gin.H{"event": "connected", "username": username})

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage && string(payload) == "ping" {
			_ = conn.WriteJSON(gin.H{"event": "pong"})
		}
	}
}
