package cache

import (
	"fmt"
	"quick-im-demo/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	conf := config.ReadConfig()
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Username: conf.Redis.Username,
		Password: conf.Redis.Password,
	})
}

func UnreadKey(username string) string { return "im:unread:" + username }
