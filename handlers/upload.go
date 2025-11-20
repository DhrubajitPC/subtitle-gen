package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadHandler handles video file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UploadHandler called")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit upload size (e.g., 100MB)
	r.ParseMultipartForm(100 << 20)

	file, header, err := r.FormFile("videoFile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory if not exists
	uploadDir := "./static/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Save file
	// Sanitize filename: replace non-alphanumeric characters (except . and -) with _
	safeFilename := func(name string) string {
		var result []rune
		for _, r := range name {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' {
				result = append(result, r)
			} else {
				result = append(result, '_')
			}
		}
		return string(result)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), safeFilename(header.Filename))
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// Render the player fragment
	tmplPath := filepath.Join("templates", "player.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"VideoPath": "/static/uploads/" + filename,
		"LocalPath": filePath, // Hidden field for backend processing
	}

	tmpl.Execute(w, data)
}
