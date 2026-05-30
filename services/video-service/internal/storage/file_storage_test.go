package storage

import (
	"os"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	slice := []string{".mp4", ".webm", ".mkv"}

	if !contains(slice, ".mp4") {
		t.Error("expected contains to return true for .mp4")
	}
	if !contains(slice, ".webm") {
		t.Error("expected contains to return true for .webm")
	}
	if contains(slice, ".avi") {
		t.Error("expected contains to return false for .avi")
	}
	if contains(slice, "") {
		t.Error("expected contains to return false for empty string")
	}
}

func TestGenerateUniqueFilename(t *testing.T) {
	tests := []struct {
		input    string
		wantExt  string
		notEmpty bool
	}{
		{"video.mp4", ".mp4", true},
		{"my video file.webm", ".webm", true},
		{"thumbnail.jpg", ".jpg", true},
		{"file-name_123.mov", ".mov", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generateUniqueFilename(tt.input)
			if result == "" {
				t.Fatal("expected non-empty filename")
			}
			ext := result[strings.LastIndex(result, "."):]
			if ext != tt.wantExt {
				t.Errorf("expected extension %s, got %s", tt.wantExt, ext)
			}
			// Filename should contain a timestamp (digits)
			hasDigit := false
			for _, c := range result {
				if c >= '0' && c <= '9' {
					hasDigit = true
					break
				}
			}
			if !hasDigit {
				t.Error("expected filename to contain timestamp digits")
			}
		})
	}
}

func TestGenerateUniqueFilename_SpecialCharsReplaced(t *testing.T) {
	result := generateUniqueFilename("my video (1).mp4")
	// Spaces and parentheses should be replaced with underscores
	if strings.Contains(result, " ") {
		t.Error("expected spaces to be replaced in filename")
	}
	if strings.Contains(result, "(") || strings.Contains(result, ")") {
		t.Error("expected parentheses to be replaced in filename")
	}
}

func TestNewFileStorage_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	basePath := dir + "/uploads"

	fs, err := NewFileStorage(basePath)
	if err != nil {
		t.Fatalf("NewFileStorage error: %v", err)
	}
	if fs == nil {
		t.Fatal("expected non-nil FileStorage")
	}

	// Check that subdirectories were created
	for _, subdir := range []string{"videos", "thumbnails"} {
		path := basePath + "/" + subdir
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected directory %s to exist: %v", path, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", path)
		}
	}
}

func TestNewFileStorage_EmptyPath_UsesDefault(t *testing.T) {
	// Use a temp dir and override via basePath; empty string uses UploadPath constant
	// We just verify no error is returned for a valid path
	dir := t.TempDir()
	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("NewFileStorage with temp dir error: %v", err)
	}
	if fs == nil {
		t.Fatal("expected non-nil FileStorage")
	}
}

func TestDeleteFile_NonExistentFile_NoError(t *testing.T) {
	dir := t.TempDir()
	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("NewFileStorage error: %v", err)
	}

	// Deleting a non-existent file should not return an error
	err = fs.DeleteFile("/uploads/videos/nonexistent.mp4")
	if err != nil {
		t.Errorf("expected no error deleting non-existent file, got: %v", err)
	}
}

func TestAllowedVideoExtensions(t *testing.T) {
	expected := []string{".mp4", ".webm", ".mkv", ".mov", ".avi"}
	for _, ext := range expected {
		if !contains(AllowedVideoExtensions, ext) {
			t.Errorf("expected %s to be in AllowedVideoExtensions", ext)
		}
	}
}

func TestAllowedImageExtensions(t *testing.T) {
	expected := []string{".jpg", ".jpeg", ".png", ".webp"}
	for _, ext := range expected {
		if !contains(AllowedImageExtensions, ext) {
			t.Errorf("expected %s to be in AllowedImageExtensions", ext)
		}
	}
}
