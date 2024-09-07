package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

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

type ByFrequencyThenRecency []FileEntry

func (a ByFrequencyThenRecency) Len() int      { return len(a) }
func (a ByFrequencyThenRecency) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByFrequencyThenRecency) Less(i, j int) bool {
	// Sort by frequency (ascending), then by last accessed (ascending)
	if a[i].Frequency != a[j].Frequency {
		return a[i].Frequency < a[j].Frequency
	}
	return a[i].LastAccessed < (a[j].LastAccessed)
}

func sortEntries(entries []FileEntry) []FileEntry {
	sort.Sort(ByFrequencyThenRecency(entries))
	return entries
}

func execute(command string) {
	entries, err := readFileStore()
	if err != nil {
		log.Fatal(err)
	}

	entries = sortEntries(entries) // Sort by frequency or use another sorting function

	if len(entries) > 0 {
		topEntry := entries[0]
		// Execute the specified command on the top entry
		logger.Log.Printf("executing command: %s %s", command, topEntry.Path)
		cmd := exec.Command(command, topEntry.Path)

		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No entries found.")
	}
}
