package main

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
	"time"
)

var dataFile string

// Struct to hold the file metadata
type FileEntry struct {
	Path         string
	Frequency    int
	LastAccessed int64
}

// Sorting Methods

type ByFrequencyThenRecency []FileEntry

func (a ByFrequencyThenRecency) Len() int      { return len(a) }
func (a ByFrequencyThenRecency) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByFrequencyThenRecency) Less(i, j int) bool {
	// Sort by frequency (descending), then by last accessed (descending)
	if a[i].Frequency != a[j].Frequency {
		return a[i].Frequency > a[j].Frequency
	}
	return a[i].LastAccessed > (a[j].LastAccessed)
}

func sortEntries(entries []FileEntry) []FileEntry {
	sort.Sort(ByFrequencyThenRecency(entries))
	return entries
}

func openTopChoice(command string) {
	entries, err := readFromStore()
	if err != nil {
		log.Fatal(err)
	}

	entries = sortEntries(entries) // Sort by frequency or use another sorting function

	if len(entries) > 0 {
		topEntry := entries[0]
		// Execute the specified command on the top entry
		cmd := exec.Command(command, topEntry.Path)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No entries found.")
	}
}

func logFileAccess(path string) {
	entries, err := readFromStore()
	if err != nil {
		log.Fatal(err)
	}

	found := false
	for i, entry := range entries {
		if entry.Path == path {
			entries[i].Frequency++
			entries[i].LastAccessed = time.Now().Unix()
			found = true
			break
		}
	}

	if !found {
		// Add a new entry if the file hasn't been logged before
		newEntry := FileEntry{
			Path:         path,
			Frequency:    1,
			LastAccessed: time.Now().Unix(),
		}
		entries = append(entries, newEntry)
	}

	// Write updated entries back to the file
	writeToStore(entries)
}
