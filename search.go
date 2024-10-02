package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

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
