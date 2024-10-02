package main

import (
	"flag"
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
