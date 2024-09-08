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

	// Internal flags
	add := flag.String("add", "", "Internal: Add path to the store")
	sanitize := flag.Bool("sanitize", false, "Internal: Sanitize command before processing")
	proc := flag.Bool("proc", false, "Internal: Process a zsh-hook command")

	// User flags
	init := flag.Bool("init", false, "Initialize fasder. Flags: zsh-hook aliases zsh-aliases, or auto for all  ")
	execCmd := flag.String("e", "", "Execute provided command against best match")
	list := flag.Bool("l", false, "List only. Omit rankings")
	reverse := flag.Bool("r", false, "Reverse sort. Useful to pipe into fzf")

	filesOnly := flag.Bool("f", false, "Files only")
	dirsOnly := flag.Bool("d", false, "Dirs only")

	flag.Parse()

	files := *filesOnly || (!*filesOnly && !*dirsOnly)
	dirs := *dirsOnly || (!*filesOnly && !*dirsOnly)

	LoadFileStore()

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
	matchingEntries := fuzzyFind(entries, searchTerm, files, dirs)

	sortedEntries := sortEntries(matchingEntries, *reverse)

	if *execCmd != "" {
		// Open the top choice with the specified command if -exec is set
		execute(sortedEntries, *execCmd)
		// logFileAccess("/path/to/file1")
	} else {
		logger.Log.Printf("displaying... %v", filesOnly)
		displaySortedEntries(sortedEntries, *list)
	}
}
