package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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
