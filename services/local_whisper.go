package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var execLookPath = exec.LookPath

// TranscribeAudioLocal uses the local 'whisper' CLI tool to transcribe audio.
func TranscribeAudioLocal(audioPath string) (string, error) {
	// Check if whisper is installed
	whisperCmd := "whisper"
	if _, err := execLookPath(whisperCmd); err != nil {
		// Fallback to common pip install location on Mac
		fallbackPath := "/Users/dhch/Library/Python/3.9/bin/whisper"
		if _, err := os.Stat(fallbackPath); err == nil {
			whisperCmd = fallbackPath
		} else {
			return "", fmt.Errorf("whisper CLI tool not found in PATH or %s. Please ensure 'openai-whisper' is installed via pip", fallbackPath)
		}
	}

	// Create a temporary directory for output
	tempDir, err := os.MkdirTemp("", "whisper_output")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Construct command
	// whisper <audioPath> --model base --output_format txt --output_dir <tempDir>
	cmd := execCommand(whisperCmd, audioPath, "--model", "base", "--output_format", "txt", "--output_dir", tempDir)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper command failed: %v\nOutput: %s", err, string(output))
	}

	// Read the output file
	// Whisper creates a file with the same basename as the audio file but with .txt extension
	baseName := filepath.Base(audioPath)
	// Remove extension from baseName to get the name whisper uses
	fileNameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	outputFilePath := filepath.Join(tempDir, fileNameWithoutExt+".txt")

	content, err := os.ReadFile(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read transcript file: %v", err)
	}

	return string(content), nil
}
