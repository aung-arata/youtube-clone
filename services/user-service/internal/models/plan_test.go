package models

import (
	"encoding/json"
	"testing"
)

func TestPlan_JSONSerialization(t *testing.T) {
	p := Plan{
		ID:                 1,
		Name:               "Premium",
		Price:              9.99,
		MaxVideoQuality:    "4K",
		MaxUploadsPerMonth: 50,
		AdsFree:            true,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded Plan
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.ID != p.ID {
		t.Errorf("ID: got %d, want %d", decoded.ID, p.ID)
	}
	if decoded.Name != p.Name {
		t.Errorf("Name: got %s, want %s", decoded.Name, p.Name)
	}
	if decoded.Price != p.Price {
		t.Errorf("Price: got %f, want %f", decoded.Price, p.Price)
	}
	if decoded.AdsFree != p.AdsFree {
		t.Errorf("AdsFree: got %v, want %v", decoded.AdsFree, p.AdsFree)
	}
	if decoded.MaxUploadsPerMonth != p.MaxUploadsPerMonth {
		t.Errorf("MaxUploadsPerMonth: got %d, want %d", decoded.MaxUploadsPerMonth, p.MaxUploadsPerMonth)
	}
}

func TestSubscription_JSONSerialization(t *testing.T) {
	s := Subscription{
		ID:          1,
		UserID:      5,
		ChannelName: "TechChannel",
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded Subscription
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.UserID != s.UserID {
		t.Errorf("UserID: got %d, want %d", decoded.UserID, s.UserID)
	}
	if decoded.ChannelName != s.ChannelName {
		t.Errorf("ChannelName: got %s, want %s", decoded.ChannelName, s.ChannelName)
	}
}
