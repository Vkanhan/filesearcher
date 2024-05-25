package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var (
	homeTemplate    *template.Template
	resultsTemplate *template.Template
)

func main() {
	// Parse the templates
	var err error
	homeTemplate, err = template.ParseFiles("templates/home.html")
	if err != nil {
		fmt.Println("Error parsing home template:", err)
		os.Exit(1)
	}

	resultsTemplate, err = template.ParseFiles("templates/results.html")
	if err != nil {
		fmt.Println("Error parsing results template:", err)
		os.Exit(1)
	}

	// Serve the HTML form at the root URL
	http.HandleFunc("/", homeHandler)
	// Handle form submission at /search URL
	http.HandleFunc("/search", searchHandler)

	// Start the web server on port 8080
	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}

// homeHandler serves the HTML form
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := homeTemplate.Execute(w, nil); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}

// searchHandler handles the form submission and performs the file search
func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	pattern := r.FormValue("pattern")
	directory := r.FormValue("directory")

	// Compile the regex pattern
	regexPattern, err := regexp.Compile(pattern)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error compiling regex pattern: %v", err), http.StatusBadRequest)
		return
	}

	// Perform the file search
	matches, err := searchFiles(regexPattern, directory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Display the search results
	data := struct {
		Matches []string
	}{
		Matches: matches,
	}

	if err := resultsTemplate.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Error displaying results: %v", err), http.StatusInternalServerError)
	}
}

// searchFiles performs the recursive file search
func searchFiles(regexPattern *regexp.Regexp, directory string) ([]string, error) {
	var matches []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file name matches the search pattern
		if !info.IsDir() && regexPattern.MatchString(info.Name()) {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}
