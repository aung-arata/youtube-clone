package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func newTestRouter(method, path string, handler http.HandlerFunc) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(path, handler).Methods(method)
	return r
}

func newVideoHandler() *VideoHandler {
	return &VideoHandler{db: nil}
}

func TestCreateVideo_BadJSON(t *testing.T) {
	h := newVideoHandler()
	req := httptest.NewRequest(http.MethodPost, "/videos", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateVideo(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %d", rr.Code)
	}
}

func TestCreateVideo_EmptyTitle(t *testing.T) {
	h := newVideoHandler()
	body, _ := json.Marshal(map[string]interface{}{
		"title":        "",
		"url":          "http://example.com/video.mp4",
		"channel_name": "TestChannel",
	})
	req := httptest.NewRequest(http.MethodPost, "/videos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateVideo(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty title, got %d", rr.Code)
	}
}

func TestCreateVideo_EmptyURL(t *testing.T) {
	h := newVideoHandler()
	body, _ := json.Marshal(map[string]interface{}{
		"title":        "My Video",
		"url":          "",
		"channel_name": "TestChannel",
	})
	req := httptest.NewRequest(http.MethodPost, "/videos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateVideo(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty url, got %d", rr.Code)
	}
}

func TestCreateVideo_EmptyChannelName(t *testing.T) {
	h := newVideoHandler()
	body, _ := json.Marshal(map[string]interface{}{
		"title":        "My Video",
		"url":          "http://example.com/video.mp4",
		"channel_name": "",
	})
	req := httptest.NewRequest(http.MethodPost, "/videos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateVideo(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty channel_name, got %d", rr.Code)
	}
}

func TestGetVideo_InvalidID(t *testing.T) {
	h := newVideoHandler()
	router := newTestRouter(http.MethodGet, "/videos/{id}", h.GetVideo)

	req := httptest.NewRequest(http.MethodGet, "/videos/notanumber", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid id, got %d", rr.Code)
	}
}

func TestIncrementViews_InvalidID(t *testing.T) {
	h := newVideoHandler()
	router := newTestRouter(http.MethodPost, "/videos/{id}/views", h.IncrementViews)

	req := httptest.NewRequest(http.MethodPost, "/videos/notanumber/views", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid id, got %d", rr.Code)
	}
}
