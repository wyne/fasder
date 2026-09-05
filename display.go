package main

import (
	"fmt"
	"time"
)

// Display

func displaySortedEntries(entries []PathEntry, list bool) {
	now := time.Now().Unix()
	for _, entry := range entries {
		if list {
			fmt.Println(entry.Path)
		} else {
			fmt.Printf(
				"%-11.5f%s\n",
				frecentScore(entry, now),
				entry.Path)
		}
	}
}
