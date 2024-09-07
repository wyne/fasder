package main

import (
	"fmt"
)

// Display

func displaySortedEntries(entries []PathEntry, list bool) {
	for _, entry := range entries {
		if list {
			fmt.Println(entry.Path)
		} else {
			fmt.Printf(
				"%.5f\t%s\n",
				entry.Rank,
				entry.Path)
		}
	}
}
