package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu    sync.RWMutex
	users map[string]map[*websocket.Conn]struct{} //使用嵌套map作为字段，将用户名与对应websocket连接和连接内容进行绑定
}

// 创建新的用户
func NewHub() *Hub {
	return &Hub{
		users: make(map[string]map[*websocket.Conn]struct{}),
	}
}

// 添加新的连接，如果当前没有则进行创建
func (h *Hub) Add(username string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.users[username] == nil {
		h.users[username] = make(map[*websocket.Conn]struct{})
	}
	h.users[username][conn] = struct{}{}
}

func (h *Hub) Remove(username string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	//删除连接只需要从用户map对应的键Conn删除即可，若删除后没有连接，则继续删除对应用户
	delete(h.users[username], conn)
	if len(h.users[username]) == 0 {
		delete(h.users, username)
	}
}

func (h *Hub) IsOnline(username string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.users[username]) > 0
}
