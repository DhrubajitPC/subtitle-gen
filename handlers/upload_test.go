package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadHandler(t *testing.T) {
	// Setup temp dir
	tmpDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup templates
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}
	// Dummy player.html
	playerContent := `<div>Video: {{.VideoPath}}</div>`
	if err := os.WriteFile(filepath.Join(templatesDir, "player.html"), []byte(playerContent), 0644); err != nil {
		t.Fatalf("Failed to write player.html: %v", err)
	}

	// Change CWD
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change wd: %v", err)
	}
	defer os.Chdir(originalWd)

	// Prepare multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("videoFile", "test_video.mp4")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("dummy video content"))
	writer.Close()

	req, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UploadHandler)

	handler.ServeHTTP(rr, req)

	// Check status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check body
	if !strings.Contains(rr.Body.String(), "Video: /static/uploads/") {
		t.Errorf("handler returned unexpected body: %v", rr.Body.String())
	}

	// Verify file exists in static/uploads
	// The handler creates static/uploads in CWD (which is tmpDir)
	// Filename is timestamp_safeFilename
	// We just check if the directory has files
	uploadsDir := filepath.Join(tmpDir, "static", "uploads")
	entries, err := os.ReadDir(uploadsDir)
	if err != nil {
		t.Fatalf("Failed to read uploads dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("No file uploaded to static/uploads")
	}
}
