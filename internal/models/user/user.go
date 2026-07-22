package user

import "time"

type User struct {
	Id        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	LoginTime *time.Time `json:"login_time"`
	CreatedAt time.Time  `json:"created_at"`
}

func (User) TableName() string {
	return "user"
}
