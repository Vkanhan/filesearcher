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
	filenameFlag := flag.String("file", "", "File name to search for")
	directoryFlag := flag.String("directory", "", "Directory to search in")
	cliMode := flag.Bool("cli", false, "Enable console mode")

	flag.Parse()

	switch {
	case *cliMode:
		runCLIMode(filenameFlag, directoryFlag)
	default:
		runWebMode()
	}
}

func runCLIMode(filenameFlag, directoryFlag *string) {
	filename, directory := getCLIInput(filenameFlag, directoryFlag)

	matches, err := searchFiles(filename, directory)
	if err != nil {
		log.Fatalf("Error searching directory: %v\n", err)
	}

	displayCLIResult(matches)
}

func displayCLIResult(matches []string) {
	if len(matches) == 0 {
		fmt.Println("No matches found")
	} else {
		for _, match := range matches {
			fmt.Println(match)
		}
	}
}

func getCLIInput(filenameFlag, directoryFlag *string) (string, string) {
	reader := bufio.NewReader(os.Stdin)

	filename := *filenameFlag
	if filename == "" {
		fmt.Print("Enter the filename to search for: ")
		filename = readCLIInput(reader)
	}

	directory := *directoryFlag
	if directory == "" {
		fmt.Print("Enter the directory to search for: ")
		directory = readCLIInput(reader)
	}

	if filename == "" || directory == "" {
		log.Fatalf("Error: Both file name and directory are required.")
	}

	return filename, directory
}

func readCLIInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading string: %v", err)
	}
	return strings.TrimSpace(input)
}

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

	matches, err := searchFiles(filename, directory)
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

func searchFiles(filename, directory string) ([]string, error) {
	var matches []string
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Skipping %s: %v\n", path, err)
				return nil // Continue the walk without stopping
			}
			return err
		}

		if !d.IsDir() && d.Name() == filename {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}
