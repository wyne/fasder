package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/sahilm/fuzzy"
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
		var pathEndSegments []string
		var pathFirstSegments []string

		// Split search term by spaces
		searchTerms := strings.Split(searchTerm, " ")

		// Prepare paths and segments
		for _, entry := range entries {
			parts := strings.Split(entry.Path, "/")
			pathEndSegments = append(pathEndSegments, parts[len(parts)-1])                         // Last segment
			pathFirstSegments = append(pathFirstSegments, strings.Join(parts[:len(parts)-1], "/")) // First part of path
		}

		var matches fuzzy.Matches
		if len(searchTerms) > 1 {
			searchTermStart := strings.Join(searchTerms[:len(searchTerms)-1], "")
			searchTermEnd := searchTerms[len(searchTerms)-1]

			// Fuzzy match first term in first part of path and second term in last segment
			leadingTermMatches := fuzzy.Find(searchTermStart, pathFirstSegments) // First part matching
			endTermMatches := fuzzy.Find(searchTermEnd, pathEndSegments)         // Last segment matching

			// Set results to paths that match both leadingTermMatches and endTermMatches
			for _, leadingMatch := range leadingTermMatches {
				for _, endMatch := range endTermMatches {
					if leadingMatch.Index == endMatch.Index { // Both match the same path
						results = append(results, entries[leadingMatch.Index])
					}
				}
			}

		} else {
			// Perform fuzzy search on the last segment only (default behavior)
			matches = fuzzy.Find(searchTerms[0], pathEndSegments)
			for _, match := range matches {
				results = append(results, entries[match.Index])
			}
		}
	}

	// Filter results based on file/dir-only flags
	return filterEntries(results, filesOnly, dirsOnly)
}

// Helper function to filter files or directories
func filterEntries(entries []PathEntry, files bool, dirs bool) []PathEntry {
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
