package main

import (
	"flag"
)

var (
	fastwalkFlag = flag.Bool("fastwalk", false, "Enable fastwalk lib")
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
