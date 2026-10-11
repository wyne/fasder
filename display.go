package main

import (
	"fmt"
	"time"
)

func displayEntries(entries []PathEntry, list bool, mode sortMode) {
	displayEntriesAt(entries, list, mode, time.Now().Unix())
}

func displayEntriesAt(entries []PathEntry, list bool, mode sortMode, now int64) {
	for _, entry := range entries {
		if list {
			fmt.Println(entry.Path)
		} else {
			score := frecentScore(entry, now)
			switch mode {
			case sortRank:
				score = entry.Rank
			case sortRecent:
				score = recentScore(entry, now)
			}
			fmt.Printf(
				"%-11.5f%s\n",
				score,
				entry.Path)
		}
	}
}
