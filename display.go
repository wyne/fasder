package main

import (
	"fmt"
	"time"
)

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
