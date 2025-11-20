package services

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExtractAudio extracts audio from a video file and saves it as an MP3.
// Returns the path to the generated audio file.
func ExtractAudio(videoPath string) (string, error) {
	// Construct output path (replace extension with .mp3)
	ext := filepath.Ext(videoPath)
	audioPath := strings.TrimSuffix(videoPath, ext) + ".mp3"

	// ffmpeg command: -i input -q:a 0 -map a output.mp3
	// -y to overwrite if exists
	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-q:a", "0", "-map", "a", audioPath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}

	return audioPath, nil
}
