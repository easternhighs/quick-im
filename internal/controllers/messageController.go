package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"quick-im-demo/internal/cache"
	"quick-im-demo/internal/models/message"
	"quick-im-demo/internal/models/user"
	"quick-im-demo/internal/ws"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MessageController struct {
	DB      *gorm.DB
	Hub     *ws.Hub
	Redis   *redis.Client
	Channel string
}

type SendMessageRequest struct {
	ToUser  string `json:"to_user" binding:"required,max=32"`
	Content string `json:"content" binding:"required,max=2000"`
}

type HistoryQuery struct {
	Peer   string `form:"peer" binding:"required,max=32"`
	Before uint64 `form:"before"`
}

type MarkReadRequest struct {
	Peer string `json:"peer" binding:"required,max=32"`
}

type ConversationItem struct {
	Peer        string    `json:"peer"`
	LastContent string    `json:"last_content"`
	LastTime    time.Time `json:"last_time"`
	Unread      int64     `json:"unread"`
}

func (mc *MessageController) Send(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromUser, toUser, content := c.GetString("username"), strings.TrimSpace(req.ToUser), strings.TrimSpace(req.Content)
	if fromUser == toUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot send a message to yourself"})
		return
	}

	var receiver user.User
	if err := mc.DB.Where("username = ?", toUser).First(&receiver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "receiver does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query receiver failed"})
		return
	}

	msgID, err := newMessageID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "generate message id failed"})
		return
	}
	msg := &message.Message{MsgID: msgID, FromUser: fromUser, ToUser: toUser, Content: content, CreatedAt: time.Now()}
	if err := mc.DB.Create(msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save message failed"})
		return
	}

	// MySQL remains the reliable source of truth; Redis enriches delivery.
	if err := mc.Redis.HIncrBy(c.Request.Context(), cache.UnreadKey(toUser), fromUser, 1).Err(); err != nil {
		log.Printf("increment unread count failed: %v", err)
	}
	body, err := json.Marshal(msg)
	if err != nil || mc.Redis.Publish(c.Request.Context(), mc.Channel, body).Err() != nil {
		log.Printf("publish message failed; using local push")
		mc.Hub.Push(msg.ToUser, "message.new", msg)
	}
	c.JSON(http.StatusCreated, gin.H{"data": msg})
}

// History uses an ID cursor so later pages are stable while new messages arrive.
func (mc *MessageController) History(c *gin.Context) {
	var query HistoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username, items := c.GetString("username"), make([]message.Message, 0)
	dbQuery := mc.DB.Where("(from_user = ? AND to_user = ?) OR (from_user = ? AND to_user = ?)", username, query.Peer, query.Peer, username).Order("id DESC").Limit(30)
	if query.Before > 0 {
		dbQuery = dbQuery.Where("id < ?", query.Before)
	}
	if err := dbQuery.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query history failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (mc *MessageController) MarkRead(c *gin.Context) {
	var req MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := mc.Redis.HDel(c.Request.Context(), cache.UnreadKey(c.GetString("username")), req.Peer).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "clear unread count failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (mc *MessageController) Unread(c *gin.Context) {
	values, err := mc.Redis.HGetAll(c.Request.Context(), cache.UnreadKey(c.GetString("username"))).Result()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "query unread count failed"})
		return
	}
	counts := make(map[string]int64, len(values))
	for peer, value := range values {
		counts[peer], _ = strconv.ParseInt(value, 10, 64)
	}
	c.JSON(http.StatusOK, gin.H{"data": counts})
}

// Conversations returns the newest message for each peer. MySQL 8 window
// functions keep this query small and make it suitable for the demo inbox.
func (mc *MessageController) Conversations(c *gin.Context) {
	username := c.GetString("username")
	items := make([]ConversationItem, 0)
	err := mc.DB.Raw(`
		SELECT peer, content AS last_content, created_at AS last_time
		FROM (
			SELECT CASE WHEN from_user = ? THEN to_user ELSE from_user END AS peer,
				content, created_at,
				ROW_NUMBER() OVER (
					PARTITION BY CASE WHEN from_user = ? THEN to_user ELSE from_user END
					ORDER BY id DESC
				) AS row_num
			FROM im_messages
			WHERE from_user = ? OR to_user = ?
		) AS ranked
		WHERE row_num = 1
		ORDER BY last_time DESC`, username, username, username, username).Scan(&items).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query conversations failed"})
		return
	}

	unread, err := mc.Redis.HGetAll(c.Request.Context(), cache.UnreadKey(username)).Result()
	if err == nil {
		for index := range items {
			items[index].Unread, _ = strconv.ParseInt(unread[items[index].Peer], 10, 64)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func newMessageID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
