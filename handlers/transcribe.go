package handlers

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"video-subtitle-generator/services"
)

func validateVideoPath(videoPath string) error {
	// 1. Basic empty check
	if videoPath == "" {
		return fmt.Errorf("video path is empty")
	}

	// 2. Prevent Argument Injection (starting with -)
	if strings.HasPrefix(videoPath, "-") {
		return fmt.Errorf("invalid filename")
	}

	// 3. Path Traversal Check
	// Clean the path to resolve .. and .
	cleanPath := filepath.Clean(videoPath)

	// Get absolute path of the uploads directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("server error resolving paths")
	}
	uploadsDir := filepath.Join(cwd, "static", "uploads")

	// Resolve symlinks for uploadsDir to ensure we have the canonical path
	absUploadsDir, err := filepath.EvalSymlinks(uploadsDir)
	if err != nil {
		// If the directory doesn't exist yet (might happen if no uploads yet),
		// we can try Abs, but for security it's better if it exists.
		// In this app, upload.go creates it.
		// If it doesn't exist, we can't validate properly, so fail safe.
		absUploadsDir, err = filepath.Abs(uploadsDir)
		if err != nil {
			return fmt.Errorf("server error resolving uploads path")
		}
	}

	// Get absolute path of the requested video
	absVideoPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("invalid video path")
	}

	// Resolve symlinks for the video path as well
	// This handles cases where the file path might contain symlinks
	realVideoPath, err := filepath.EvalSymlinks(absVideoPath)
	if err != nil {
		// If file doesn't exist, EvalSymlinks fails.
		// We check existence later, but for prefix check we need canonical path.
		// If it fails, it likely doesn't exist or is invalid.
		return fmt.Errorf("file not found or invalid path")
	}

	// Check if realVideoPath starts with absUploadsDir
	if !strings.HasPrefix(realVideoPath, absUploadsDir) {
		return fmt.Errorf("access denied: file must be in uploads directory")
	}

	// Check if file actually exists (redundant if EvalSymlinks succeeded, but good for clarity)
	if _, err := os.Stat(realVideoPath); err != nil {
		return fmt.Errorf("file not found")
	}

	return nil
}

func TranscribeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("TranscribeHandler called")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	videoPath := r.FormValue("videoPath") // This should be the local filesystem path

	// Security Validation
	if err := validateVideoPath(videoPath); err != nil {
		log.Printf("Security alert: Invalid video path attempt: %s. Error: %v", videoPath, err)
		// Escape the error message to prevent XSS
		escapedErr := html.EscapeString(err.Error())
		w.Write([]byte("<div class='error'>Error: " + escapedErr + "</div>"))
		return
	}

	// 1. Extract Audio
	audioPath, err := services.ExtractAudio(videoPath)
	if err != nil {
		escapedErr := html.EscapeString(err.Error())
		w.Write([]byte("<div class='error'>Error extracting audio: " + escapedErr + "</div>"))
		return
	}

	// 2. Transcribe (Local)
	text, err := services.TranscribeAudioLocal(audioPath)
	if err != nil {
		escapedErr := html.EscapeString(err.Error())
		w.Write([]byte("<div class='error'>Error transcribing: " + escapedErr + "</div>"))
		return
	}

	// 3. Render Transcript
	tmplPath := filepath.Join("templates", "transcript.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		w.Write([]byte("<div class='error'>Template error</div>"))
		return
	}

	data := map[string]string{
		"Transcript": text,
	}
	tmpl.Execute(w, data)
}
