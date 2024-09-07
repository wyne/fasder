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

var dataFile string

// Struct to hold the file metadata
type FileEntry struct {
	Path         string
	Frequency    int
	LastAccessed int64 // Unix timestamp
}

// Sorting Methods

type ByFrequencyThenRecency struct {
	entries []FileEntry
	reverse bool
}

func (a ByFrequencyThenRecency) Len() int { return len(a.entries) }
func (a ByFrequencyThenRecency) Swap(i, j int) {
	a.entries[i], a.entries[j] = a.entries[j], a.entries[i]
}

func (a ByFrequencyThenRecency) Less(i, j int) bool {
	if a.reverse {
		// Sort in descending order
		if a.entries[i].Frequency != a.entries[j].Frequency {
			return a.entries[i].Frequency > a.entries[j].Frequency
		}
		return a.entries[i].LastAccessed > a.entries[j].LastAccessed
	} else {
		// Sort in ascending order
		if a.entries[i].Frequency != a.entries[j].Frequency {
			return a.entries[i].Frequency < a.entries[j].Frequency
		}
		return a.entries[i].LastAccessed < a.entries[j].LastAccessed
	}
}

func sortEntries(entries []FileEntry, reverse bool) []FileEntry {
	sort.Sort(ByFrequencyThenRecency{entries, reverse})
	return entries
}

// Fuzzy search function
func fuzzyFind(entries []FileEntry, searchTerm string, filesOnly bool, dirsOnly bool) []FileEntry {
	// Collect matching entries
	var results []FileEntry

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
func filterEntries(entries []FileEntry, files bool, dirs bool) []FileEntry {
	logger.Log.Printf("Filtering files: %v, dirs: %v", files, dirs)
	var filtered []FileEntry

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

func execute(entries []FileEntry, command string) {
	if len(entries) > 0 {
		bestMatch := entries[len(entries)-1]
		// Execute the specified command on the top entry
		cmdStr := fmt.Sprintf("%s %s", command, bestMatch.Path)
		parts := strings.Split(cmdStr, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No entries found.")
	}
}
