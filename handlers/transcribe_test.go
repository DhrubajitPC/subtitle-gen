package handlers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateVideoPath(t *testing.T) {
	// Setup: Create a temporary directory structure mimicking the app
	// We need a "static/uploads" directory in the current working directory for the test
	// But we don't want to mess up the actual project.
	// So we will change the CWD for the test.

	tmpDir, err := os.MkdirTemp("", "test_app")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create static/uploads inside tmpDir
	uploadsDir := filepath.Join(tmpDir, "static", "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		t.Fatalf("Failed to create uploads dir: %v", err)
	}

	// Create a valid video file
	validVideoPath := filepath.Join(uploadsDir, "valid_video.mp4")
	if err := os.WriteFile(validVideoPath, []byte("dummy content"), 0644); err != nil {
		t.Fatalf("Failed to create dummy video file: %v", err)
	}

	// Create a file outside uploads
	outsideFile := filepath.Join(tmpDir, "outside.txt")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0644); err != nil {
		t.Fatalf("Failed to create outside file: %v", err)
	}

	// Change CWD to tmpDir so the handler logic works as expected
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change wd: %v", err)
	}
	defer os.Chdir(originalWd)

	tests := []struct {
		name      string
		videoPath string
		wantErr   bool
	}{
		{
			name:      "Valid path",
			videoPath: "./static/uploads/valid_video.mp4",
			wantErr:   false,
		},
		{
			name:      "Valid absolute path",
			videoPath: validVideoPath,
			wantErr:   false,
		},
		{
			name:      "Path traversal attempt",
			videoPath: "./static/uploads/../../outside.txt",
			wantErr:   true,
		},
		{
			name:      "Argument injection attempt",
			videoPath: "-option",
			wantErr:   true,
		},
		{
			name:      "Non-existent file",
			videoPath: "./static/uploads/non_existent.mp4",
			wantErr:   true,
		},
		{
			name:      "Empty path",
			videoPath: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVideoPath(tt.videoPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateVideoPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
