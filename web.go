package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Vkanhan/filesearcher/internal/search"
)

func runWebMode() {
	const port = "8080"
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/search", fileSearchHandler)

	server := &http.Server{
		Addr: ":" + port,
	}

	fmt.Println("Starting the server at port: ", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	const homeTemplate = "templates/home.html"
	renderTemplate(w, homeTemplate, nil)
}

func fileSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	filename := r.FormValue("filename")
	directory := r.FormValue("directory")

	if !filepath.IsAbs(directory) {
		http.Error(w, "Invalid directory path", http.StatusBadRequest)
		return
	}

	matches, err := search.SearchFiles(filename, directory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching directory: %v", err), http.StatusInternalServerError)
		return
	}

	const resultTemplate = "templates/results.html"
	renderTemplate(w, resultTemplate, struct{ Matches []string }{Matches: matches})
}

func renderTemplate(w http.ResponseWriter, templatePath string, data any) {
	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing the template: %v", err), http.StatusInternalServerError)
		return
	}
	if err := templ.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("error rendering the template: %v", err), http.StatusInternalServerError)
	}
}
