package models

import (
	"encoding/json"
	"testing"
)

func TestNotification_JSONSerialization(t *testing.T) {
	n := Notification{
		ID:      1,
		UserID:  5,
		Type:    "comment",
		Title:   "New comment",
		Message: "Someone commented on your video",
		Link:    "/video/123",
		IsRead:  false,
	}

	data, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded Notification
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.ID != n.ID {
		t.Errorf("ID: got %d, want %d", decoded.ID, n.ID)
	}
	if decoded.Type != n.Type {
		t.Errorf("Type: got %s, want %s", decoded.Type, n.Type)
	}
	if decoded.Title != n.Title {
		t.Errorf("Title: got %s, want %s", decoded.Title, n.Title)
	}
	if decoded.IsRead != n.IsRead {
		t.Errorf("IsRead: got %v, want %v", decoded.IsRead, n.IsRead)
	}
}

func TestNotification_LinkOmittedWhenEmpty(t *testing.T) {
	n := Notification{
		ID:      2,
		UserID:  3,
		Type:    "like",
		Title:   "New like",
		Message: "Someone liked your video",
	}

	data, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	// link field should be omitted when empty
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if _, ok := raw["link"]; ok {
		t.Error("expected 'link' to be omitted when empty")
	}
}

func TestCreateNotificationRequest_JSONSerialization(t *testing.T) {
	req := CreateNotificationRequest{
		UserID:  10,
		Type:    "subscription",
		Title:   "New subscriber",
		Message: "You have a new subscriber",
		Link:    "/channel/me",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded CreateNotificationRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.UserID != req.UserID {
		t.Errorf("UserID: got %d, want %d", decoded.UserID, req.UserID)
	}
	if decoded.Type != req.Type {
		t.Errorf("Type: got %s, want %s", decoded.Type, req.Type)
	}
	if decoded.Link != req.Link {
		t.Errorf("Link: got %s, want %s", decoded.Link, req.Link)
	}
}
