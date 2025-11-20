package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"video-subtitle-generator/handlers"
)

func main() {
	// Define the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/transcribe", handlers.TranscribeHandler)

	// Start server
	fmt.Printf("Server starting on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
