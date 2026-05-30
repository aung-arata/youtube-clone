package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestVideo_JSONSerialization(t *testing.T) {
	userID := 7
	v := Video{
		ID:               1,
		UserID:           &userID,
		Title:            "Test Video",
		Description:      "A test description",
		URL:              "/uploads/videos/test.mp4",
		Thumbnail:        "/uploads/thumbnails/test.jpg",
		ChannelName:      "TestChannel",
		Visibility:       "public",
		ProcessingStatus: "ready",
		Views:            100,
		Likes:            10,
		Dislikes:         1,
		Category:         "education",
		Duration:         "5:30",
		UploadedAt:       time.Now().UTC(),
	}

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded Video
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.ID != v.ID {
		t.Errorf("ID: got %d, want %d", decoded.ID, v.ID)
	}
	if decoded.Title != v.Title {
		t.Errorf("Title: got %s, want %s", decoded.Title, v.Title)
	}
	if decoded.Views != v.Views {
		t.Errorf("Views: got %d, want %d", decoded.Views, v.Views)
	}
	if decoded.ProcessingStatus != v.ProcessingStatus {
		t.Errorf("ProcessingStatus: got %s, want %s", decoded.ProcessingStatus, v.ProcessingStatus)
	}
}

func TestVideo_UserIDOmittedWhenNil(t *testing.T) {
	v := Video{
		ID:    2,
		Title: "No user video",
	}

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if _, ok := raw["user_id"]; ok {
		t.Error("expected 'user_id' to be omitted when nil")
	}
}

func TestPlaylist_JSONSerialization(t *testing.T) {
	p := Playlist{
		ID:          1,
		UserID:      5,
		Name:        "My Playlist",
		Description: "A great playlist",
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded Playlist
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.ID != p.ID {
		t.Errorf("ID: got %d, want %d", decoded.ID, p.ID)
	}
	if decoded.Name != p.Name {
		t.Errorf("Name: got %s, want %s", decoded.Name, p.Name)
	}
}

func TestPlaylistVideo_JSONSerialization(t *testing.T) {
	pv := PlaylistVideo{
		ID:         1,
		PlaylistID: 10,
		VideoID:    20,
		Position:   1,
	}

	data, err := json.Marshal(pv)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded PlaylistVideo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.PlaylistID != pv.PlaylistID {
		t.Errorf("PlaylistID: got %d, want %d", decoded.PlaylistID, pv.PlaylistID)
	}
	if decoded.Position != pv.Position {
		t.Errorf("Position: got %d, want %d", decoded.Position, pv.Position)
	}
}
