package main

import (
	"fmt"
)

// Display

func displaySortedEntries(entries []FileEntry) {
	entries = sortEntries(entries)
	for _, entry := range entries {
		fmt.Printf(
			"%d\t%s\n",
			entry.Frequency,
			entry.Path)
	}
}
