package integration_test

import (
"bytes"
"context"
"mime/multipart"
"net/http"
"net/http/httptest"
"os"
"testing"

"github.com/aung-arata/youtube-clone/backend/internal/handlers"
"github.com/aung-arata/youtube-clone/backend/internal/middleware"
"github.com/aung-arata/youtube-clone/backend/internal/storage"
"github.com/DATA-DOG/go-sqlmock"
)

func TestUploadIntegration_UploadVideo(t *testing.T) {
// Create mock database
db, mock, err := sqlmock.New()
if err != nil {
t.Fatalf("Failed to create mock database: %v", err)
}
defer db.Close()

// Create temporary directory for file storage
tmpDir := "/tmp/test_uploads"
os.RemoveAll(tmpDir)
defer os.RemoveAll(tmpDir)

fileStorage, err := storage.NewFileStorage(tmpDir)
if err != nil {
t.Fatalf("Failed to create file storage: %v", err)
}

handler := handlers.NewUploadHandler(db, fileStorage)

t.Run("Successful Video Upload", func(t *testing.T) {
// Create multipart form data
body := &bytes.Buffer{}
writer := multipart.NewWriter(body)

// Add form fields
writer.WriteField("title", "Test Video")
writer.WriteField("description", "Test Description")
writer.WriteField("category", "Technology")
writer.WriteField("channel_name", "Test Channel")
writer.WriteField("channel_avatar", "https://example.com/avatar.jpg")
writer.WriteField("duration", "10:30")

// Add fake video file
videoWriter, _ := writer.CreateFormFile("video", "test.mp4")
videoWriter.Write([]byte("fake video content"))

// Add fake thumbnail
thumbnailWriter, _ := writer.CreateFormFile("thumbnail", "thumbnail.jpg")
thumbnailWriter.Write([]byte("fake image content"))

writer.Close()

// Mock database insert
mock.ExpectQuery("INSERT INTO videos").
WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "url", "thumbnail", "channel_name", "channel_avatar", "views", "likes", "dislikes", "category", "duration", "uploaded_at", "created_at", "updated_at"}).
AddRow(1, "Test Video", "Test Description", "/uploads/videos/test.mp4", "/uploads/thumbnails/thumbnail.jpg", "Test Channel", "https://example.com/avatar.jpg", 0, 0, 0, "Technology", "10:30", "2024-01-01", "2024-01-01", "2024-01-01"))

req := httptest.NewRequest(http.MethodPost, "/upload/video", body)
req.Header.Set("Content-Type", writer.FormDataContentType())

// Add user context (simulating AuthMiddleware)
ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
req = req.WithContext(ctx)

w := httptest.NewRecorder()

handler.UploadVideo(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
}
})

t.Run("Upload without authentication", func(t *testing.T) {
body := &bytes.Buffer{}
writer := multipart.NewWriter(body)
writer.WriteField("title", "Test Video")
writer.Close()

req := httptest.NewRequest(http.MethodPost, "/upload/video", body)
req.Header.Set("Content-Type", writer.FormDataContentType())

w := httptest.NewRecorder()

handler.UploadVideo(w, req)

if w.Code != http.StatusUnauthorized {
t.Errorf("Expected status 401, got %d", w.Code)
}
})

t.Run("Upload without required fields", func(t *testing.T) {
body := &bytes.Buffer{}
writer := multipart.NewWriter(body)
// Missing title and other required fields
writer.Close()

req := httptest.NewRequest(http.MethodPost, "/upload/video", body)
req.Header.Set("Content-Type", writer.FormDataContentType())

// Add user context
ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
req = req.WithContext(ctx)

w := httptest.NewRecorder()

handler.UploadVideo(w, req)

if w.Code != http.StatusBadRequest {
t.Errorf("Expected status 400, got %d", w.Code)
}
})
}

func TestUploadIntegration_FileValidation(t *testing.T) {
tmpDir := "/tmp/test_uploads_validation"
os.RemoveAll(tmpDir)
defer os.RemoveAll(tmpDir)

fileStorage, err := storage.NewFileStorage(tmpDir)
if err != nil {
t.Fatalf("Failed to create file storage: %v", err)
}

t.Run("Invalid video file extension", func(t *testing.T) {
body := &bytes.Buffer{}
writer := multipart.NewWriter(body)

// Try to upload a file with invalid extension
fileWriter, _ := writer.CreateFormFile("video", "test.txt")
fileWriter.Write([]byte("not a video"))
writer.Close()

// Get the file from the multipart form
req := httptest.NewRequest(http.MethodPost, "/upload", body)
req.Header.Set("Content-Type", writer.FormDataContentType())
req.ParseMultipartForm(10 << 20)

file, header, _ := req.FormFile("video")
defer file.Close()

_, err := fileStorage.SaveVideo(file, header)
if err == nil {
t.Error("Expected error for invalid video extension")
}
})

t.Run("Invalid image file extension", func(t *testing.T) {
body := &bytes.Buffer{}
writer := multipart.NewWriter(body)

fileWriter, _ := writer.CreateFormFile("thumbnail", "test.txt")
fileWriter.Write([]byte("not an image"))
writer.Close()

req := httptest.NewRequest(http.MethodPost, "/upload", body)
req.Header.Set("Content-Type", writer.FormDataContentType())
req.ParseMultipartForm(10 << 20)

file, header, _ := req.FormFile("thumbnail")
defer file.Close()

_, err := fileStorage.SaveThumbnail(file, header)
if err == nil {
t.Error("Expected error for invalid image extension")
}
})
}
