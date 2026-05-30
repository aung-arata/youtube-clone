package websocket

import (
	"testing"
	"time"
)

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

	// Register client and wait for confirmation via a sync message on broadcast
	hub.register <- client

	// Poll until the hub has processed the registration
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		hub.mu.RLock()
		_, ok := hub.clients[client.UserID]
		hub.mu.RUnlock()
		if ok {
			break
		}
		time.Sleep(time.Millisecond)
	}

	hub.mu.RLock()
	clients, ok := hub.clients[client.UserID]
	hub.mu.RUnlock()
	if !ok {
		t.Fatal("expected client to be registered")
	}
	if !clients[client] {
		t.Error("expected client to be in the user's client set")
	}

	// Unregister client; hub closes client.Send when done
	hub.unregister <- client

	// Wait for Send channel to be closed — that's the hub's signal that it processed unregister
	select {
	case _, open := <-client.Send:
		if open {
			// drained a message; channel isn't closed yet — keep waiting
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for hub to unregister client")
	}

	hub.mu.RLock()
	_, still := hub.clients[client.UserID]
	hub.mu.RUnlock()
	if still {
		t.Error("expected client to be unregistered")
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

	// Wait for registration
	for i := 0; i < 1000; i++ {
		hub.mu.RLock()
		_, ok := hub.clients[client.UserID]
		hub.mu.RUnlock()
		if ok {
			break
		}
	}

	msg := []byte(`{"type":"notification","payload":"hello"}`)
	hub.broadcast <- &BroadcastMessage{UserID: 2, Message: msg}

	// Read message from client's send channel
	received := <-client.Send
	if string(received) != string(msg) {
		t.Errorf("received %q, want %q", received, msg)
	}
}

func TestHub_BroadcastToUnregisteredUser_NoError(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Sending to a user with no registered clients should not panic or block
	hub.broadcast <- &BroadcastMessage{UserID: 9999, Message: []byte("msg")}
	// If we reach here, no panic occurred
}
