package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Define flags for CLI mode
	patternFlag := flag.String("pattern", "", "Regex pattern to search for")
	directoryFlag := flag.String("directory", "", "Directory to search in")
	cliMode := flag.Bool("cli", false, "Enable console mode")

	flag.Parse()

	if *cliMode {
		// If CLI mode is enabled, run the search from the console
		runCLIMode(patternFlag, directoryFlag)
	} else {
		// If not in CLI mode, start the web server
		http.HandleFunc("/", homeHandler)
		http.HandleFunc("/search", searchHandler)

		fmt.Println("Starting server at :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Failed to start server:", err)
			os.Exit(1)
		}
	}
}

// runCLIMode performs the file search and outputs the results to the console
func runCLIMode(patternFlag, directoryFlag *string) {
	// Get user input if flags are not provided
	pattern := *patternFlag
	directory := *directoryFlag

	reader := bufio.NewReader(os.Stdin)

	if pattern == "" {
		fmt.Print("Enter the regex pattern to search for: ")
		patternInput, _ := reader.ReadString('\n')
		pattern = strings.TrimSpace(patternInput)
	}

	if directory == "" {
		fmt.Print("Enter the directory to search in: ")
		directoryInput, _ := reader.ReadString('\n')
		directory = strings.TrimSpace(directoryInput)
	}

	// Validate the input
	if pattern == "" || directory == "" {
		fmt.Println("Error: Both pattern and directory are required.")
		os.Exit(1)
	}

	// Compile the regex pattern
	regexPattern, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error compiling regex pattern: %v\n", err)
		os.Exit(1)
	}

	// Perform the file search
	matches, err := searchFiles(regexPattern, directory)
	if err != nil {
		fmt.Printf("Error searching directory: %v\n", err)
		os.Exit(1)
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
	pattern := r.FormValue("pattern")
	directory := r.FormValue("directory")

	// Validate and sanitize directory path
	if !filepath.IsAbs(directory) {
		http.Error(w, "Invalid directory path", http.StatusBadRequest)
		return
	}

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
func searchFiles(regexPattern *regexp.Regexp, directory string) ([]string, error) {
	var matches []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
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
		if !info.IsDir() && regexPattern.MatchString(info.Name()) {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}
