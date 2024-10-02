package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define flags for CLI mode
	filenameFlag := flag.String("file", "", "File name to search for")
	directoryFlag := flag.String("directory", "", "Directory to search in")
	cliMode := flag.Bool("cli", false, "Enable console mode")

	flag.Parse()

	if *cliMode {
		runCLIMode(filenameFlag, directoryFlag)
	} else {
		http.HandleFunc("/", homeHandler)
		http.HandleFunc("/search", searchHandler)

		fmt.Println("Starting server at :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("failed to start the server")
		}
	}
}

// runCLIMode performs the file search and outputs the results to the console
func runCLIMode(filenameFlag, directoryFlag *string) error {
	// Get user input if flags are not provided
	filename := *filenameFlag
	directory := *directoryFlag

	reader := bufio.NewReader(os.Stdin)

	if filename == "" {
		fmt.Print("Enter the file name to search for: ")
		filenameInput, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading the filename input: %w", err)
		}
		filename = strings.TrimSpace(filenameInput)
	}

	if directory == "" {
		fmt.Print("Enter the directory to search in: ")
		directoryInput, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading directory input: %w", err)
		}
		directory = strings.TrimSpace(directoryInput)
	}

	// Validate the input
	if filename == "" || directory == "" {
		log.Fatalf("Error: Both file name and directory are required.")
	}

	// Perform the file search
	matches, err := searchFiles(filename, directory)
	if err != nil {
		log.Fatalf("Error searching directory: %v\n", err)
	}

	// Display the search results
	if len(matches) == 0 {
		fmt.Println("No matches found.")
	} else {
		fmt.Println("Matches found:")
		for _, match := range matches {
			fmt.Println(match)
		}
	}
	return nil
}

// homeHandler serves the HTML form
func homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing home template: %v", err), http.StatusInternalServerError)
		return
	}
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
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}
	filename := r.FormValue("filename")
	directory := r.FormValue("directory")

	// Validate and sanitize directory path
	if !filepath.IsAbs(directory) {
		http.Error(w, "Invalid directory path", http.StatusBadRequest)
		return
	}

	// Perform the file search
	matches, err := searchFiles(filename, directory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse the results template
	resultsTemplate, err := template.ParseFiles("templates/results.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing results template: %v", err), http.StatusInternalServerError)
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
func searchFiles(filename, directory string) ([]string, error) {
	var matches []string
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip directories/files that cannot be accessed
			if os.IsPermission(err) {
				fmt.Printf("Skipping %s: %v\n", path, err)
				return nil // Continue the walk without stopping
			}
			// Return the error for any other issues
			return err
		}

		// Check if the file name matches the search pattern
		if !d.IsDir() && d.Name() == filename {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}
