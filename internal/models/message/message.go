package message

import "time"

type Message struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MsgID     string    `gorm:"size:36;uniqueIndex;not null" json:"msg_id"`
	FromUser  string    `gorm:"size:32;index;not null" json:"from_user"`
	ToUser    string    `gorm:"size:32;index;not null" json:"to_user"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (Message) TableName() string {
	return "im_messages"
}
