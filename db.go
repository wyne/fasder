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
		pathSegments := splitPath(entry.Path)

		// Check if search terms match the path segments in order
		if fuzzyMatchInOrder(searchTerms, pathSegments) {
			results = append(results, entry)
		}
	}

	return results
}

// Split the path into segments using both `/` and `.` as delimiters
func splitPath(path string) []string {
	// Split the path by `/` to get directory segments
	dirSegments := strings.Split(path, "/")

	var allSegments []string
	for i, dirSegment := range dirSegments {
		// For the last segment (usually a file), keep it intact (do not split by `.`)
		if i == len(dirSegments)-1 {
			allSegments = append(allSegments, dirSegment)
		} else {
			// Otherwise, split by `.` for regular segments (to handle directories like `.config`)
			fileSegments := strings.Split(dirSegment, ".")
			allSegments = append(allSegments, fileSegments...)
		}
	}

	return allSegments
}

// Check if the search terms match the path segments in order
func fuzzyMatchInOrder(searchTerms []string, pathSegments []string) bool {
	searchIndex := 0
	pathLength := len(pathSegments)

	if pathLength == 0 {
		return false
	}

	// Match each search term with the corresponding path segment
	for i := 0; i < len(searchTerms)-1; i++ {
		if searchIndex >= len(pathSegments) {
			return false
		}

		// Use fuzzy.Find to check for fuzzy matches between the current search term and the current path segment
		matches := fuzzy.Find(searchTerms[i], []string{pathSegments[searchIndex]})
		if len(matches) > 0 {
			searchIndex++
		}
	}

	// Ensure that the last search term matches the last or second-to-last path segment
	lastSearchTerm := searchTerms[len(searchTerms)-1]
	if searchIndex < len(pathSegments) {
		// Last segment match
		lastSegmentMatches := fuzzy.Find(lastSearchTerm, []string{pathSegments[len(pathSegments)-1]})
		if len(lastSegmentMatches) > 0 {
			return true
		}
	}

	if len(pathSegments) > 1 && searchIndex <= len(pathSegments)-2 {
		// Second-to-last segment match if available
		secondLastSegmentMatches := fuzzy.Find(lastSearchTerm, []string{pathSegments[len(pathSegments)-2]})
		if len(secondLastSegmentMatches) > 0 {
			return true
		}
	}

	// If all terms matched but last term didn’t match the end segments correctly
	return false
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
