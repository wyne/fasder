package main

import (
	"fmt"
)

// Display

func displaySortedEntries(entries []FileEntry) {
	entries = sortEntries(entries)
	for _, entry := range entries {
		fmt.Printf("Path: %s, Frequency: %d, Last Accessed: %d\n",
			entry.Path, entry.Frequency, entry.LastAccessed)
	}
}
