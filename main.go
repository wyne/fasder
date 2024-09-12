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
	version := flag.Bool("v", false, "View version")
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

	if *version {
		println("0.1.3")
		return
	}

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

	// Read from file store
	entries, err := readFileStore()
	if err != nil {
		log.Fatal(err)
	}

	// Search
	searchTerm := strings.Join(flag.Args(), " ")
	logger.Log.Printf("Search term: %s", searchTerm)
	matchingEntries := fuzzyFind(entries, searchTerm)
	filteredEntries := filterEntries(matchingEntries, files, dirs)
	sortedEntries := sortEntries(filteredEntries, *reverse)

	// Execute if necessary
	if *execCmd != "" {
		execute(sortedEntries, *execCmd)
	} else {
		displaySortedEntries(sortedEntries, *list)
	}
}
