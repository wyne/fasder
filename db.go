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

// Fuzzy find function that matches the search terms to the paths
func fuzzyFind(entries []PathEntry, searchTerm string) []PathEntry {
	// Collect matching entries
	var results []PathEntry

	if searchTerm == "" {
		return entries
	}

	// Split the search term into segments based on spaces
	searchTerms := strings.Split(searchTerm, " ")

	// Loop through each entry and check for a match
	for _, entry := range entries {
		// Split the path into segments using `/` for directories but treat the last segment (file) as a whole
		pathSegments := splitByPath(entry.Path)
		pathAndFileSegments := splitByPathAndFile(entry.Path)

		methodA := matchInOrder(searchTerms, pathSegments)
		methodB := matchInOrder(searchTerms, pathAndFileSegments)

		// Check if search terms match the path segments in order
		if methodA || methodB {
			results = append(results, entry)
		}
	}

	return results
}

// Split the path into segments using '/'
func splitByPath(path string) []string {
	return strings.Split(path, "/")
}

// Split the path into segments using both `/` and `.` as delimiters
// Split by both '/' and '.'
func splitByPathAndFile(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '.'
	})
}

// Check if the search terms match the path segments in order,
// with the last search term matching the last path segment.
func matchInOrder(searchTerms []string, pathSegments []string) bool {
	// Ensure both searchTerms and pathSegments are not empty
	if len(searchTerms) == 0 || len(pathSegments) == 0 {
		return false
	}

	// Ensure the last search term matches the last path segment
	if !match(searchTerms[len(searchTerms)-1], pathSegments[len(pathSegments)-1]) {
		return false
	}

	// Initialize index pointers for search terms and path segments
	searchIndex := 0
	segmentIndex := 0

	// Match the search terms to path segments in order (ignoring adjacency but maintaining sequence)
	for searchIndex < len(searchTerms)-1 && segmentIndex < len(pathSegments)-1 {
		if match(searchTerms[searchIndex], pathSegments[segmentIndex]) {
			// Move to the next search term if a match is found
			searchIndex++
		}
		// Always move to the next path segment
		segmentIndex++
	}

	// If all search terms were matched, return true
	return searchIndex == len(searchTerms)-1
}

// Fuzzy match function using fuzzy.Find
func match(searchTerm string, pathSegment string) bool {
	// Perform fuzzy finding on the path segment using the search term
	matches := fuzzy.Find(searchTerm, []string{pathSegment})
	return len(matches) > 0 // Return true if there's a match, otherwise false
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
