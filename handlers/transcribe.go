package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"video-subtitle-generator/services"
)

func TranscribeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("TranscribeHandler called")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	videoPath := r.FormValue("videoPath") // This should be the local filesystem path

	// 1. Extract Audio
	audioPath, err := services.ExtractAudio(videoPath)
	if err != nil {
		w.Write([]byte("<div class='error'>Error extracting audio: " + err.Error() + "</div>"))
		return
	}

	// 2. Transcribe (Local)
	text, err := services.TranscribeAudioLocal(audioPath)
	if err != nil {
		w.Write([]byte("<div class='error'>Error transcribing: " + err.Error() + "</div>"))
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
