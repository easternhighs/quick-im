package user

type UserMessage struct {
	Content     string `json:"content"`
	ContentType int    `json:"content_type"`
	FromUserID  int64  `json:"from_user_id"`
	ToUserID    int64  `json:"to_user_id"`
	DeliverTime int64  `json:"deliver_time"`
	IsRevoked   bool   `json:"is_revoked"`
}
