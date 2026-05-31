package websocket

import (
	"testing"
	"time"
)

// waitForCondition polls fn up to deadline, sleeping 1ms between attempts.
func waitForCondition(t *testing.T, deadline time.Duration, fn func() bool) bool {
	t.Helper()
	end := time.Now().Add(deadline)
	for time.Now().Before(end) {
		if fn() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("expected non-nil Hub")
	}
	if hub.clients == nil {
		t.Error("expected clients map to be initialized")
	}
	if hub.broadcast == nil {
		t.Error("expected broadcast channel to be initialized")
	}
	if hub.register == nil {
		t.Error("expected register channel to be initialized")
	}
	if hub.unregister == nil {
		t.Error("expected unregister channel to be initialized")
	}
}

func TestHub_RegisterAndUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:      hub,
		Send:     make(chan []byte, 10),
		UserID:   1,
		ClientID: "test-client-1",
	}

	hub.register <- client

	if !waitForCondition(t, 2*time.Second, func() bool {
		hub.mu.RLock()
		defer hub.mu.RUnlock()
		return hub.clients[client.UserID] != nil
	}) {
		t.Fatal("timed out waiting for client to be registered")
	}

	hub.mu.RLock()
	clients := hub.clients[client.UserID]
	hub.mu.RUnlock()
	if !clients[client] {
		t.Error("expected client to be in the user's client set")
	}

	// Unregister; hub closes client.Send as its acknowledgement
	hub.unregister <- client

	select {
	case _, open := <-client.Send:
		// channel closed (open==false) means hub processed the unregister
		if open {
			// unexpected message; drain and wait for close
			select {
			case <-client.Send:
			case <-time.After(2 * time.Second):
				t.Fatal("timed out waiting for Send channel to close")
			}
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for hub to unregister client")
	}

	hub.mu.RLock()
	_, still := hub.clients[client.UserID]
	hub.mu.RUnlock()
	if still {
		t.Error("expected client to be removed from hub after unregister")
	}
}

func TestHub_BroadcastToRegisteredUser(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:      hub,
		Send:     make(chan []byte, 10),
		UserID:   2,
		ClientID: "test-client-2",
	}

	hub.register <- client

	if !waitForCondition(t, 2*time.Second, func() bool {
		hub.mu.RLock()
		defer hub.mu.RUnlock()
		return hub.clients[client.UserID] != nil
	}) {
		t.Fatal("timed out waiting for client to be registered")
	}

	msg := []byte(`{"type":"notification","payload":"hello"}`)
	hub.broadcast <- &BroadcastMessage{UserID: 2, Message: msg}

	select {
	case received := <-client.Send:
		if string(received) != string(msg) {
			t.Errorf("received %q, want %q", received, msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting to receive broadcast message")
	}
}

func TestHub_BroadcastToUnregisteredUser_NoError(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Sending to a user with no registered clients should not panic or block
	hub.broadcast <- &BroadcastMessage{UserID: 9999, Message: []byte("msg")}
	// If we reach here, no panic occurred
}
