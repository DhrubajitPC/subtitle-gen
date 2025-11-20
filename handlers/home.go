package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	layoutPath := filepath.Join("templates", "layout.html")

	tmpl, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", nil)
	if err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
