package services

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// TestHelperProcess isn't a real test. It's used to mock exec.Command.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Check if we are mocking ffmpeg
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd := args[0]
	if cmd == "ffmpeg" {
		// Verify arguments if needed
		// Just exit success
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
	os.Exit(2)
}

func TestExtractAudio(t *testing.T) {
	// Mock execCommand
	execCommand = func(name string, arg ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
	defer func() { execCommand = exec.Command }()

	// Test
	videoPath := "test_video.mp4"
	audioPath, err := ExtractAudio(videoPath)
	if err != nil {
		t.Errorf("ExtractAudio failed: %v", err)
	}
	if audioPath != "test_video.mp3" {
		t.Errorf("Expected audio path 'test_video.mp3', got '%s'", audioPath)
	}
}
