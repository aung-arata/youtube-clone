package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWatchHistory_JSONSerialization(t *testing.T) {
	wh := WatchHistory{
		ID:        1,
		UserID:    5,
		VideoID:   10,
		WatchedAt: time.Now().UTC().Truncate(time.Second),
	}

	data, err := json.Marshal(wh)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded WatchHistory
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.ID != wh.ID {
		t.Errorf("ID: got %d, want %d", decoded.ID, wh.ID)
	}
	if decoded.UserID != wh.UserID {
		t.Errorf("UserID: got %d, want %d", decoded.UserID, wh.UserID)
	}
	if decoded.VideoID != wh.VideoID {
		t.Errorf("VideoID: got %d, want %d", decoded.VideoID, wh.VideoID)
	}
}

func TestVideoWithHistory_JSONSerialization(t *testing.T) {
	vwh := VideoWithHistory{
		Video: Video{
			ID:          1,
			Title:       "History Video",
			ChannelName: "TestChannel",
			Views:       500,
			Duration:    "10:00",
		},
		WatchedAt: time.Now().UTC().Truncate(time.Second),
	}

	data, err := json.Marshal(vwh)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded VideoWithHistory
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.Title != vwh.Title {
		t.Errorf("Title: got %s, want %s", decoded.Title, vwh.Title)
	}
	if decoded.Views != vwh.Views {
		t.Errorf("Views: got %d, want %d", decoded.Views, vwh.Views)
	}
}
