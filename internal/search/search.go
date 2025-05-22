package search

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/charlievieth/fastwalk"
)

type WalkerFunc func(root string, fn func(path string, d fs.DirEntry, err error) error) error

func FastWalker(root string, fn func(path string, d fs.DirEntry, err error) error) error {
	conf := &fastwalk.Config{Follow: false}
	return fastwalk.Walk(conf, root, fn)
}

func StdlibWalker(root string, fn func(path string, d fs.DirEntry, err error) error) error {
	return filepath.WalkDir(root, fn)
}

func SearchFiles(filename, directory string, walker WalkerFunc) ([]string, error) {
	var matches []string
	err := walker(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
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
