package storage

import (
"fmt"
"io"
"mime/multipart"
"os"
"path/filepath"
"strings"
"time"
)

const (
// MaxVideoSize is the maximum allowed video file size (500 MB)
MaxVideoSize = 500 * 1024 * 1024
// MaxThumbnailSize is the maximum allowed thumbnail size (5 MB)
MaxThumbnailSize = 5 * 1024 * 1024
// UploadPath is the base path for uploaded files
UploadPath = "./uploads"
)

var (
// AllowedVideoExtensions defines allowed video file extensions
AllowedVideoExtensions = []string{".mp4", ".webm", ".mkv", ".mov", ".avi"}
// AllowedImageExtensions defines allowed image file extensions
AllowedImageExtensions = []string{".jpg", ".jpeg", ".png", ".webp"}
)

// FileStorage handles file upload and storage operations
type FileStorage struct {
basePath string
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(basePath string) (*FileStorage, error) {
if basePath == "" {
basePath = UploadPath
}

// Create base directory if it doesn't exist
if err := os.MkdirAll(basePath, 0755); err != nil {
return nil, fmt.Errorf("failed to create upload directory: %w", err)
}

// Create subdirectories
videosPath := filepath.Join(basePath, "videos")
thumbnailsPath := filepath.Join(basePath, "thumbnails")

if err := os.MkdirAll(videosPath, 0755); err != nil {
return nil, fmt.Errorf("failed to create videos directory: %w", err)
}

if err := os.MkdirAll(thumbnailsPath, 0755); err != nil {
return nil, fmt.Errorf("failed to create thumbnails directory: %w", err)
}

return &FileStorage{basePath: basePath}, nil
}

// SaveVideo saves an uploaded video file
func (fs *FileStorage) SaveVideo(file multipart.File, header *multipart.FileHeader) (string, error) {
// Validate file size
if header.Size > MaxVideoSize {
return "", fmt.Errorf("video file too large: max size is %d MB", MaxVideoSize/(1024*1024))
}

// Validate file extension
ext := strings.ToLower(filepath.Ext(header.Filename))
if !contains(AllowedVideoExtensions, ext) {
return "", fmt.Errorf("invalid video format: allowed formats are %v", AllowedVideoExtensions)
}

// Generate unique filename
filename := generateUniqueFilename(header.Filename)
filePath := filepath.Join(fs.basePath, "videos", filename)

// Create the file
dst, err := os.Create(filePath)
if err != nil {
return "", fmt.Errorf("failed to create file: %w", err)
}
defer dst.Close()

// Copy file content
if _, err := io.Copy(dst, file); err != nil {
return "", fmt.Errorf("failed to save file: %w", err)
}

// Return the relative URL
return "/uploads/videos/" + filename, nil
}

// SaveThumbnail saves an uploaded thumbnail image
func (fs *FileStorage) SaveThumbnail(file multipart.File, header *multipart.FileHeader) (string, error) {
// Validate file size
if header.Size > MaxThumbnailSize {
return "", fmt.Errorf("thumbnail file too large: max size is %d MB", MaxThumbnailSize/(1024*1024))
}

// Validate file extension
ext := strings.ToLower(filepath.Ext(header.Filename))
if !contains(AllowedImageExtensions, ext) {
return "", fmt.Errorf("invalid image format: allowed formats are %v", AllowedImageExtensions)
}

// Generate unique filename
filename := generateUniqueFilename(header.Filename)
filePath := filepath.Join(fs.basePath, "thumbnails", filename)

// Create the file
dst, err := os.Create(filePath)
if err != nil {
return "", fmt.Errorf("failed to create file: %w", err)
}
defer dst.Close()

// Copy file content
if _, err := io.Copy(dst, file); err != nil {
return "", fmt.Errorf("failed to save file: %w", err)
}

// Return the relative URL
return "/uploads/thumbnails/" + filename, nil
}

// DeleteFile deletes a file from storage
func (fs *FileStorage) DeleteFile(url string) error {
// Convert URL to file path
filePath := filepath.Join(fs.basePath, strings.TrimPrefix(url, "/uploads/"))

if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
return fmt.Errorf("failed to delete file: %w", err)
}

return nil
}

// generateUniqueFilename generates a unique filename with timestamp
func generateUniqueFilename(originalName string) string {
ext := filepath.Ext(originalName)
name := strings.TrimSuffix(originalName, ext)
// Clean the filename
name = strings.Map(func(r rune) rune {
if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
return r
}
return '_'
}, name)

timestamp := time.Now().Unix()
return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
for _, s := range slice {
if s == item {
return true
}
}
return false
}
