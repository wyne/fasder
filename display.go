package main

import (
	"fmt"
	"time"
)

// Display

func displaySortedEntries(entries []PathEntry, list bool) {
	displayEntries(entries, list, false)
}

func displayRankSortedEntries(entries []PathEntry, list bool) {
	displayEntries(entries, list, true)
}

func displayEntries(entries []PathEntry, list bool, rankScore bool) {
	now := time.Now().Unix()
	for _, entry := range entries {
		if list {
			fmt.Println(entry.Path)
		} else {
			score := frecentScore(entry, now)
			if rankScore {
				score = entry.Rank
			}
			fmt.Printf(
				"%-11.5f%s\n",
				score,
				entry.Path)
		}
	}
}
