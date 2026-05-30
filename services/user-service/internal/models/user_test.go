package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUser_PasswordNotInJSON(t *testing.T) {
	u := User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: "supersecret",
		Role:     "user",
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if _, ok := raw["password"]; ok {
		t.Error("password field must not appear in JSON output (json:\"-\" tag expected)")
	}
}

func TestUser_JSONSerialization(t *testing.T) {
	planID := 2
	u := User{
		ID:        10,
		Username:  "bob",
		Email:     "bob@example.com",
		Password:  "hidden",
		Avatar:    "/avatars/bob.png",
		Role:      "admin",
		PlanID:    &planID,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if decoded.ID != u.ID {
		t.Errorf("ID: got %d, want %d", decoded.ID, u.ID)
	}
	if decoded.Username != u.Username {
		t.Errorf("Username: got %s, want %s", decoded.Username, u.Username)
	}
	if decoded.Role != u.Role {
		t.Errorf("Role: got %s, want %s", decoded.Role, u.Role)
	}
	if decoded.PlanID == nil || *decoded.PlanID != planID {
		t.Errorf("PlanID: got %v, want %d", decoded.PlanID, planID)
	}
}

func TestUser_PlanIDOmittedWhenNil(t *testing.T) {
	u := User{
		ID:       3,
		Username: "noPlan",
		Email:    "noplan@example.com",
		Role:     "user",
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if _, ok := raw["plan_id"]; ok {
		t.Error("expected 'plan_id' to be omitted when nil")
	}
}
