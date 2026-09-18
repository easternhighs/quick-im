package ws

import "testing"

func TestHubTracksOnlineUser(t *testing.T) {
	hub := NewHub()
	session := NewSession(nil)

	hub.Add("alice", session)
	if !hub.IsOnline("alice") {
		t.Fatal("alice should be online after Add")
	}

	hub.Remove("alice", session)
	if hub.IsOnline("alice") {
		t.Fatal("alice should be offline after Remove")
	}
}

func TestSessionEnqueueStopsAfterClose(t *testing.T) {
	session := NewSession(nil)
	if !session.Enqueue(map[string]string{"event": "test"}) {
		t.Fatal("enqueue should succeed")
	}
	session.Close()
	if session.Enqueue(map[string]string{"event": "test"}) {
		t.Fatal("closed session must reject messages")
	}
}
