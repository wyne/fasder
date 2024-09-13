package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wyne/fasder/logger"
)

var dataFile string

func LoadFileStore() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Silently return
		return
	}

	// Expand the ~ to the home directory
	dataFile = filepath.Join(homeDir, ".fasder")
}

// Reads the `.fasder` file and loads file entries into a slice

func readFileStore() ([]PathEntry, error) {
	var entries []PathEntry
	f, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil // File doesn't exist yet, return empty list
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue // Skip malformed lines
		}

		freq, _ := strconv.ParseFloat(parts[1], 64)
		lastAccessed, _ := strconv.ParseInt(parts[2], 10, 64)

		entry := PathEntry{
			Path:         parts[0],
			Rank:         freq,
			LastAccessed: lastAccessed,
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

var mu sync.Mutex

func writeFileStore(entries []PathEntry) {
	mu.Lock() // Lock to prevent concurrent access
	defer mu.Unlock()

	tempPrefix := "fasder-"

	var cumulativeRank float64
	for _, entry := range entries {
		cumulativeRank += entry.Rank
	}

	// Apply decay if cumulative rank exceeds threshold
	const threshold = 2000.0
	const decayFactor = 0.9

	if cumulativeRank > threshold {
		logger.Log.Println("Rank threshold met. Decaying...")
		for i := range entries {
			entries[i].Rank *= decayFactor
		}
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp(filepath.Dir(dataFile), tempPrefix)
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()

	for _, entry := range entries {
		line := fmt.Sprintf("%s|%.5f|%d\n", entry.Path, entry.Rank, entry.LastAccessed)
		if _, err := tempFile.WriteString(line); err != nil {
			log.Fatal(err)
		}
	}

	// Sync to make sure all data is written
	if err := tempFile.Sync(); err != nil {
		log.Fatal(err)
	}

	// Close the temporary file before renaming
	if err := tempFile.Close(); err != nil {
		log.Fatal(err)
	}

	// Rename the temporary file to replace the original file atomically
	if err := os.Rename(tempFile.Name(), dataFile); err != nil {
		log.Fatal(err)
	}
}

// AddToStore an entry to the store
func AddToStore(path string) {
	entries, err := readFileStore()
	if err != nil {
		log.Fatal(err)
	}

	found := false
	for i, entry := range entries {
		if entry.Path == path {
			logger.Log.Printf(
				"Adding path: %s %v->%v",
				path,
				entries[i].Rank,
				entries[i].Rank+1/entries[i].Rank,
			)
			entries[i].Rank = entries[i].Rank + 1/entries[i].Rank

			entries[i].LastAccessed = time.Now().Unix()
			found = true
			break
		}
	}

	if !found {
		// Add a new entry if the file hasn't been logged before
		newEntry := PathEntry{
			Path:         path,
			Rank:         1,
			LastAccessed: time.Now().Unix(),
		}
		entries = append(entries, newEntry)
	}

	// Write updated entries back to the file
	writeFileStore(entries)
}
