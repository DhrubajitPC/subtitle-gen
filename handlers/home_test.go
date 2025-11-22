package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// Setup temp dir with templates
	tmpDir, err := os.MkdirTemp("", "home_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	// Create dummy layout.html
	layoutContent := `{{define "layout.html"}}<html><body>{{template "content" .}}</body></html>{{end}}`
	if err := os.WriteFile(filepath.Join(templatesDir, "layout.html"), []byte(layoutContent), 0644); err != nil {
		t.Fatalf("Failed to write layout.html: %v", err)
	}

	// Create dummy index.html
	indexContent := `{{define "content"}}<h1>Welcome</h1>{{end}}`
	if err := os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(indexContent), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}

	// Change CWD
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change wd: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HomeHandler)

	// Serve
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check body
	expected := "<h1>Welcome</h1>"
	if !contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
	// simple contains check, or use strings.Contains
}
