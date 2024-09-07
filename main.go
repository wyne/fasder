package main

import (
	"flag"
	"log"
	"strings"

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

	logger.Log.Printf("Positional arguments: %v\n", flag.Args())

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

	searchTerm := strings.Join(flag.Args(), " ")

	logger.Log.Printf("searching.... %s", searchTerm)
	matchingEntries := fuzzyFind(entries, searchTerm)

	sortedEntries := sortEntries(matchingEntries)

	if *execCmd != "" {
		// Open the top choice with the specified command if -exec is set
		execute(sortedEntries, *execCmd)
		// logFileAccess("/path/to/file1")
	} else {
		logger.Log.Printf("displaying... %v", filesOnly)
		displaySortedEntries(sortedEntries)
	}
}
