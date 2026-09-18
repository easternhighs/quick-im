package controllers

import (
	"net/http"
	"quick-im-demo/internal/utils"
	"quick-im-demo/internal/ws"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

const pongWait = 60 * time.Second

type WebSocketController struct{ Hub *ws.Hub }

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return origin == "" || origin == "http://"+r.Host || origin == "https://"+r.Host
}}

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
	session := ws.NewSession(conn)
	wc.Hub.Add(username, session)
	defer wc.Hub.Remove(username, session)
	go session.WritePump()
	session.Enqueue(gin.H{"event": "connected", "username": username})

	conn.SetReadLimit(8 * 1024)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { return conn.SetReadDeadline(time.Now().Add(pongWait)) })
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage && string(payload) == "ping" {
			session.Enqueue(gin.H{"event": "pong"})
		}
	}
}
