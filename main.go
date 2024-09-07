package main

import (
	"flag"
	"log"

	"github.com/wyne/fasder/logger"
)

// Global variable to hold the logger
var Logger *log.Logger

func main() {
	logger.InitLog()

	// Implement command-line flags
	add := flag.String("add", "", "Add object to store")
	init := flag.Bool("init", false, "Initializers: zsh-hook")
	sanitize := flag.Bool("sanitize", false, "Sanitize command before processing")
	proc := flag.Bool("proc", false, "Process a command")
	execCmd := flag.String("e", "", "Command to open the top choice")

	filesOnly := flag.Bool("f", false, "Files")

	LoadFileStore()

	flag.Parse()

	// Commands

	if *init {
		Init(flag.Args())
		return
	}

	if *sanitize {
		Sanitize(flag.Args())
		return
	}

	if *proc {
		Proc(flag.Args())
		return
	}

	if *add != "" {
		AddToStore(*add)
		return
	}

	// Retrieve entries

	entries, err := readFileStore()
	if err != nil {
		log.Fatal(err)
	}

	if *execCmd != "" {
		// Open the top choice with the specified command if -exec is set
		execute(*execCmd)
		// logFileAccess("/path/to/file1")
	} else {
		logger.Log.Printf("displaying... %v", filesOnly)
		displaySortedEntries(entries)
	}
}
