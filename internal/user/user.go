package user

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	LoginTime int64  `json:"login_time"`
}
