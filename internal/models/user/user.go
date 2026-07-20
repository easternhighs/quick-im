package user

import "time"

type User struct {
	Id        int64      `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	LoginTime *time.Time `json:"login_time"`
	CreatedAt time.Time  `json:"created_at"`
}

func (User) TableName() string {
	return "user"
}
