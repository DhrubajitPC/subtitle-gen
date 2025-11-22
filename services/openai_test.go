package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTranscribeAudio(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		// Verify auth header
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header, got %s", r.Header.Get("Authorization"))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"text": "Hello world"}`)
	}))
	defer ts.Close()

	// Override endpoint
	originalEndpoint := OpenAIEndpoint
	OpenAIEndpoint = ts.URL
	defer func() { OpenAIEndpoint = originalEndpoint }()

	// Create dummy audio file
	tmpFile, err := os.CreateTemp("", "audio.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write([]byte("dummy audio"))
	tmpFile.Close()

	// Call function
	text, err := TranscribeAudio(tmpFile.Name(), "test-api-key")
	if err != nil {
		t.Fatalf("TranscribeAudio failed: %v", err)
	}

	if text != "Hello world" {
		t.Errorf("Expected 'Hello world', got '%s'", text)
	}
}

func TestTranscribeAudioError(t *testing.T) {
	// Mock server returning error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error": "bad request"}`)
	}))
	defer ts.Close()

	OpenAIEndpoint = ts.URL

	tmpFile, _ := os.CreateTemp("", "audio.mp3")
	defer os.Remove(tmpFile.Name())

	_, err := TranscribeAudio(tmpFile.Name(), "key")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
