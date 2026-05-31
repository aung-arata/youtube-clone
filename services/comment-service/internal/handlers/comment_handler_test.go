package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// newTestRouter creates a mux router with the given path and handler, used to
// inject mux URL variables for handler tests.
func newTestRouter(method, path string, handler http.HandlerFunc) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(path, handler).Methods(method)
	return r
}

func newCommentHandler() *CommentHandler {
	return &CommentHandler{db: nil}
}

func TestGetComments_InvalidVideoID(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodGet, "/videos/{videoId}/comments", h.GetComments)

	req := httptest.NewRequest(http.MethodGet, "/videos/notanumber/comments", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid videoId, got %d", rr.Code)
	}
}

func TestGetComment_InvalidID(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodGet, "/comments/{id}", h.GetComment)

	req := httptest.NewRequest(http.MethodGet, "/comments/notanumber", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid id, got %d", rr.Code)
	}
}

func TestCreateComment_InvalidVideoID(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodPost, "/videos/{videoId}/comments", h.CreateComment)

	body, _ := json.Marshal(map[string]interface{}{"content": "hello", "user_id": 1})
	req := httptest.NewRequest(http.MethodPost, "/videos/bad/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid videoId, got %d", rr.Code)
	}
}

func TestCreateComment_BadJSON(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodPost, "/videos/{videoId}/comments", h.CreateComment)

	req := httptest.NewRequest(http.MethodPost, "/videos/1/comments", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %d", rr.Code)
	}
}

func TestCreateComment_EmptyContent(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodPost, "/videos/{videoId}/comments", h.CreateComment)

	body, _ := json.Marshal(map[string]interface{}{"content": "", "user_id": 1})
	req := httptest.NewRequest(http.MethodPost, "/videos/1/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty content, got %d", rr.Code)
	}
}

func TestCreateComment_MissingUserID(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodPost, "/videos/{videoId}/comments", h.CreateComment)

	body, _ := json.Marshal(map[string]interface{}{"content": "hello", "user_id": 0})
	req := httptest.NewRequest(http.MethodPost, "/videos/1/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero user_id, got %d", rr.Code)
	}
}

func TestUpdateComment_InvalidID(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodPut, "/comments/{id}", h.UpdateComment)

	body, _ := json.Marshal(map[string]interface{}{"content": "updated"})
	req := httptest.NewRequest(http.MethodPut, "/comments/notanumber", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid id, got %d", rr.Code)
	}
}

func TestUpdateComment_EmptyContent(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodPut, "/comments/{id}", h.UpdateComment)

	body, _ := json.Marshal(map[string]interface{}{"content": ""})
	req := httptest.NewRequest(http.MethodPut, "/comments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty content, got %d", rr.Code)
	}
}

func TestDeleteComment_InvalidID(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodDelete, "/comments/{id}", h.DeleteComment)

	req := httptest.NewRequest(http.MethodDelete, "/comments/notanumber", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid id, got %d", rr.Code)
	}
}

func TestGetReplies_InvalidCommentID(t *testing.T) {
	h := newCommentHandler()
	router := newTestRouter(http.MethodGet, "/comments/{commentId}/replies", h.GetReplies)

	req := httptest.NewRequest(http.MethodGet, "/comments/notanumber/replies", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid commentId, got %d", rr.Code)
	}
}
