package realtime

import (
	"context"
	"encoding/json"
	"log"
	"quick-im-demo/internal/models/message"
	"quick-im-demo/internal/ws"

	"github.com/redis/go-redis/v9"
)

// SubscribeMessages runs once per process. Each process only pushes events
// to the connections held by its own Hub.
func SubscribeMessages(ctx context.Context, rdb *redis.Client, channel string, hub *ws.Hub) {
	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()
	if _, err := sub.Receive(ctx); err != nil {
		log.Printf("redis subscription unavailable: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-sub.Channel():
			if !ok {
				return
			}
			var msg message.Message
			if err := json.Unmarshal([]byte(item.Payload), &msg); err != nil {
				log.Printf("invalid redis message: %v", err)
				continue
			}
			hub.Push(msg.ToUser, "message.new", &msg)
		}
	}
}
