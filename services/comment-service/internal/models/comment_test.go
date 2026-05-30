package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestComment_JSONSerialization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	c := Comment{
		ID:        1,
		VideoID:   10,
		UserID:    5,
		Username:  "alice",
		Content:   "Great video!",
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded Comment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.ID != c.ID {
		t.Errorf("ID: got %d, want %d", decoded.ID, c.ID)
	}
	if decoded.VideoID != c.VideoID {
		t.Errorf("VideoID: got %d, want %d", decoded.VideoID, c.VideoID)
	}
	if decoded.Username != c.Username {
		t.Errorf("Username: got %s, want %s", decoded.Username, c.Username)
	}
	if decoded.Content != c.Content {
		t.Errorf("Content: got %s, want %s", decoded.Content, c.Content)
	}
}

func TestComment_ParentID_Optional(t *testing.T) {
	// A top-level comment has no parent
	c := Comment{
		ID:      1,
		VideoID: 1,
		UserID:  1,
		Content: "top level",
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	// parent_id should be omitted when nil
	if string(data) == "" {
		t.Fatal("expected non-empty JSON")
	}

	var decoded Comment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if decoded.ParentID != nil {
		t.Error("expected ParentID to be nil for top-level comment")
	}
}

func TestComment_ParentID_SetForReply(t *testing.T) {
	parentID := 42
	c := Comment{
		ID:       2,
		VideoID:  1,
		ParentID: &parentID,
		UserID:   3,
		Content:  "reply comment",
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded Comment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if decoded.ParentID == nil {
		t.Fatal("expected ParentID to be set for reply")
	}
	if *decoded.ParentID != parentID {
		t.Errorf("ParentID: got %d, want %d", *decoded.ParentID, parentID)
	}
}
