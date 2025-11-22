package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHelperProcessWhisper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
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
	// Check if it's whisper (or the fallback path)
	if strings.Contains(cmd, "whisper") {
		// Parse args to find output dir and audio path
		// args: [whisper, audioPath, --model, base, --output_format, txt, --output_dir, tempDir]
		var outputDir string
		var audioPath string
		for i, arg := range args {
			if arg == "--output_dir" && i+1 < len(args) {
				outputDir = args[i+1]
			}
			if i == 1 { // audioPath is usually the second arg (index 1)
				audioPath = arg
			}
		}

		if outputDir != "" && audioPath != "" {
			// Create the output file
			baseName := filepath.Base(audioPath)
			fileNameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
			outputFile := filepath.Join(outputDir, fileNameWithoutExt+".txt")

			err := os.WriteFile(outputFile, []byte("Transcribed text"), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write output file: %v\n", err)
				os.Exit(1)
			}
		}
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
	os.Exit(2)
}

func TestTranscribeAudioLocal(t *testing.T) {
	// Mock execLookPath
	execLookPath = func(file string) (string, error) {
		return "/usr/bin/whisper", nil
	}
	defer func() { execLookPath = exec.LookPath }()

	// Mock execCommand
	execCommand = func(name string, arg ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcessWhisper", "--", name}
		cs = append(cs, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
	defer func() { execCommand = exec.Command }()

	// Create dummy audio file
	tmpFile, err := os.CreateTemp("", "test_audio.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Test
	text, err := TranscribeAudioLocal(tmpFile.Name())
	if err != nil {
		t.Errorf("TranscribeAudioLocal failed: %v", err)
	}
	if text != "Transcribed text" {
		t.Errorf("Expected 'Transcribed text', got '%s'", text)
	}
}
