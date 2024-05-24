package main 

import(
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	// Define the root command for the CLI application
	var rootCmd = &cobra.Command{
		Use:   "filesearch",
		Short: "Search for files matching a pattern",
		Run:   runSearch,
	}

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runSearch function performs the file search
func runSearch(cmd *cobra.Command, args []string) {
	var pattern, directory string

	// Prompt for search pattern and directory if not provided as arguments
	if len(args) < 2 {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter the search pattern: ")
		pattern, _ = reader.ReadString('\n')
		pattern = strings.TrimSpace(pattern)

		fmt.Print("Enter the directory to search: ")
		directory, _ = reader.ReadString('\n')
		directory = strings.TrimSpace(directory)
	} else {
		// Use the provided arguments as search pattern and directory
		pattern = args[0]
		directory = args[1]
	}

	// Pattern matching logic
	var regexPattern *regexp.Regexp
	var err error
	regexPattern, err = regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error compiling regex pattern: %v\n", err)
		return
	}

	// Recursive file search
	var matches []string
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		// Check if the file name matches the search pattern
		if !info.IsDir() && regexPattern.MatchString(info.Name()) {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error searching directory: %v\n", err)
		return
	}

	// Write matching paths
	if len(matches) == 0 {
		fmt.Println("No matches found.")
	} else {
		for _, match := range matches {
			fmt.Println(match)
		}
	}
}

