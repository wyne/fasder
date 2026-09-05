package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/sahilm/fuzzy"
)

// Struct to hold the file metadata
type PathEntry struct {
	Path         string
	Rank         float64
	LastAccessed int64 // Unix timestamp
}

func sortEntries(entries []PathEntry, reverse bool) []PathEntry {
	sortEntriesAt(entries, reverse, time.Now().Unix())
	return entries
}

func sortEntriesAt(entries []PathEntry, reverse bool, now int64) []PathEntry {
	sort.Sort(ByFrecentScore{entries, reverse, now})
	return entries
}

type ByFrecentScore struct {
	entries []PathEntry
	reverse bool
	now     int64
}

func (a ByFrecentScore) Len() int { return len(a.entries) }
func (a ByFrecentScore) Swap(i, j int) {
	a.entries[i], a.entries[j] = a.entries[j], a.entries[i]
}

func (a ByFrecentScore) Less(i, j int) bool {
	scoreI := frecentScore(a.entries[i], a.now)
	scoreJ := frecentScore(a.entries[j], a.now)
	if a.reverse {
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		return a.entries[i].LastAccessed > a.entries[j].LastAccessed
	}
	if scoreI != scoreJ {
		return scoreI < scoreJ
	}
	return a.entries[i].LastAccessed < a.entries[j].LastAccessed
}

func frecentScore(entry PathEntry, now int64) float64 {
	return entry.Rank * frecentMultiplier(entry.LastAccessed, now)
}

func frecentMultiplier(lastAccessed int64, now int64) float64 {
	elapsed := now - lastAccessed
	if elapsed < 3600 {
		return 6
	}
	if elapsed < 86400 {
		return 4
	}
	if elapsed < 604800 {
		return 2
	}
	return 1
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
