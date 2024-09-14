package main

import (
	"log"
	"os"
	"strings"

	flag "github.com/cornfeedhobo/pflag"
	"github.com/wyne/fasder/logger"
	"golang.org/x/term"
)

// Global variable to hold the logger
var Logger *log.Logger

func main() {
	logger.InitLog()

	// Internal flags
	add := flag.StringP("add", "A", "", "Internal: Add path to the store")
	sanitize := flag.BoolP("sanitize", "", false, "Internal: Sanitize command before processing")
	proc := flag.BoolP("proc", "", false, "Internal: Process a zsh-hook command")

	// User flags
	version := flag.BoolP("version", "v", false, "View version")
	init := flag.BoolP("init", "", false, "Initialize fasder. Args: auto aliases")
	execCmd := flag.StringP("exec", "e", "", "Execute provided command against best match")
	list := flag.BoolP("list", "l", false, "List only. Omit rankings")
	reverse := flag.BoolP("reverse", "R", false, "Reverse sort. Useful to pipe into fzf")
	scores := flag.BoolP("scores", "s", false, "Show rank scores")

	filesOnly := flag.BoolP("files", "f", false, "Files only")
	dirsOnly := flag.BoolP("directories", "d", false, "Dirs only")

	flag.Parse()

	files := *filesOnly || (!*filesOnly && !*dirsOnly)
	dirs := *dirsOnly || (!*filesOnly && !*dirsOnly)

	LoadFileStore()

	// Commands

	if *version {
		println("0.1.4")
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

	var onlyOne bool
	onlyOne = false

	// to omit score ranks in output
	if !term.IsTerminal(int(os.Stdout.Fd())) && !*list {
		onlyOne = true
	}

	// If running in a subshell (ex: vim `f zsh`), only
	// return one result, and auto apply -l list mode
	// to omit score ranks in output
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		*list = true
	}

	// Execute if necessary
	if *execCmd != "" {
		execute(sortedEntries, *execCmd)
		return
	}

	if len(sortedEntries) == 0 {
		return
	}

	if onlyOne {
		bestMatch := []PathEntry{sortedEntries[len(sortedEntries)-1]}
		displaySortedEntries(bestMatch, *list)
	} else {
		displaySortedEntries(sortedEntries, !*scores)
	}
}
