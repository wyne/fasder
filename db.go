package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/wyne/fasder/logger"
)

// Struct to hold the file metadata
type PathEntry struct {
	Path         string
	Rank         float64
	LastAccessed int64 // Unix timestamp
}

// Sorting Methods

type ByFrequencyThenRecency struct {
	entries []PathEntry
	reverse bool
}

func (a ByFrequencyThenRecency) Len() int { return len(a.entries) }
func (a ByFrequencyThenRecency) Swap(i, j int) {
	a.entries[i], a.entries[j] = a.entries[j], a.entries[i]
}

func (a ByFrequencyThenRecency) Less(i, j int) bool {
	if a.reverse {
		// Sort in descending order
		if a.entries[i].Rank != a.entries[j].Rank {
			return a.entries[i].Rank > a.entries[j].Rank
		}
		return a.entries[i].LastAccessed > a.entries[j].LastAccessed
	} else {
		// Sort in ascending order
		if a.entries[i].Rank != a.entries[j].Rank {
			return a.entries[i].Rank < a.entries[j].Rank
		}
		return a.entries[i].LastAccessed < a.entries[j].LastAccessed
	}
}

func sortEntries(entries []PathEntry, reverse bool) []PathEntry {
	sort.Sort(ByFrequencyThenRecency{entries, reverse})
	return entries
}

// Fuzzy search function
func fuzzyFind(entries []PathEntry, searchTerm string, filesOnly bool, dirsOnly bool) []PathEntry {
	// Collect matching entries
	var results []PathEntry

	if searchTerm == "" {
		results = entries
	} else {

		var paths []string
		for _, entry := range entries {
			parts := strings.Split(entry.Path, "/")
			paths = append(paths, parts[len(parts)-1])
		}

		// Perform fuzzy search
		matches := fuzzy.Find(searchTerm, paths)

		for _, match := range matches {
			results = append(results, entries[match.Index])
		}
	}

	return filterEntries(results, filesOnly, dirsOnly)
}

// Helper function to filter files or directories
func filterEntries(entries []PathEntry, files bool, dirs bool) []PathEntry {
	logger.Log.Printf("Filtering files: %v, dirs: %v", files, dirs)
	var filtered []PathEntry

	for _, entry := range entries {
		info, err := os.Stat(entry.Path)
		if err != nil {
			// Handle error (e.g., if the file does not exist)
			continue
		}

		// Filter based on the flags
		if dirs && info.IsDir() {
			filtered = append(filtered, entry)
		} else if files && !info.IsDir() {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func execute(entries []PathEntry, command string) {
	if len(entries) > 0 {
		bestMatch := entries[len(entries)-1]

		// Increment rank
		Add(bestMatch.Path)

		// Execute the specified command on the top entry
		cmdStr := fmt.Sprintf("%s %s", command, bestMatch.Path)
		parts := strings.Split(cmdStr, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No entries found.")
	}
}
