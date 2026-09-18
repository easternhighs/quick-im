package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pingPeriod = 45 * time.Second
)

// Session owns one WebSocket writer. Business goroutines only enqueue bytes,
// avoiding Gorilla WebSocket's concurrent writer restriction.
type Session struct {
	Conn *websocket.Conn
	Send chan []byte
	done chan struct{}
	once sync.Once
}

func NewSession(conn *websocket.Conn) *Session {
	return &Session{Conn: conn, Send: make(chan []byte, 64), done: make(chan struct{})}
}

func (s *Session) Enqueue(payload any) bool {
	body, err := json.Marshal(payload)
	if err != nil {
		return false
	}
	select {
	case <-s.done:
		return false
	default:
	}
	select {
	case <-s.done:
		return false
	case s.Send <- body:
		return true
	default:
		// A slow client must not block every other client.
		s.Close()
		return false
	}
}

func (s *Session) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	defer s.Close()

	for {
		select {
		case <-s.done:
			return
		case body := <-s.Send:
			_ = s.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.Conn.WriteMessage(websocket.TextMessage, body); err != nil {
				return
			}
		case <-ticker.C:
			_ = s.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Session) Close() {
	s.once.Do(func() {
		close(s.done)
		if s.Conn != nil {
			_ = s.Conn.Close()
		}
	})
}

type Hub struct {
	mu    sync.RWMutex
	users map[string]map[*Session]struct{}
}

func NewHub() *Hub { return &Hub{users: make(map[string]map[*Session]struct{})} }

func (h *Hub) Add(username string, session *Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.users[username] == nil {
		h.users[username] = make(map[*Session]struct{})
	}
	h.users[username][session] = struct{}{}
}

func (h *Hub) Remove(username string, session *Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.users[username], session)
	if len(h.users[username]) == 0 {
		delete(h.users, username)
	}
	session.Close()
}

func (h *Hub) IsOnline(username string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.users[username]) > 0
}

func (h *Hub) Push(username string, event string, payload any) {
	h.mu.RLock()
	items := make([]*Session, 0, len(h.users[username]))
	for session := range h.users[username] {
		items = append(items, session)
	}
	h.mu.RUnlock()

	for _, session := range items {
		session.Enqueue(map[string]any{"event": event, "data": payload})
	}
}
