package main

import (
	"fmt"
)

// Display

func displaySortedEntries(entries []FileEntry, list bool) {
	for _, entry := range entries {
		if list {
			fmt.Println(entry.Path)
		} else {
			fmt.Printf(
				"%d\t%s\n",
				entry.Frequency,
				entry.Path)
		}
	}
}
