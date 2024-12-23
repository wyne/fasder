package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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

	options := []string{
		os.Getenv("_FASDER_DATA"),
		filepath.Join(homeDir, ".fasder"),
		filepath.Join(os.Getenv("XDG_DATA_HOME"), "fasder", "data"),
		filepath.Join(homeDir, ".local", "share", "fasder", "data"),
	}

	for _, option := range options {
		if option != "" {
			if _, err := os.Stat(option); err == nil {
				dataFile = option
				break
			}
		}
	}

	if dataFile == "" {
		// No valid data file found, create one
		if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
			dataFile = filepath.Join(xdgDataHome, "fasder", "data")
		} else {
			dataFile = filepath.Join(homeDir, ".local", "share", "fasder", "data")
		}

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(dataFile), 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}

		// Create the file
		file, err := os.Create(dataFile)
		if err != nil {
			log.Fatalf("Failed to create data file: %v", err)
		}
		file.Close()
	}

	// Check if the file exists and is owned by the current user
	if fileInfo, err := os.Stat(dataFile); err == nil {
		if !fileInfo.Mode().IsRegular() {
			log.Fatalf("%s is not a regular file", dataFile)
		}
		currentUser, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		fileOwner, err := user.LookupId(fmt.Sprint(fileInfo.Sys().(*syscall.Stat_t).Uid))
		if err != nil {
			log.Fatal(err)
		}
		if currentUser.Uid != fileOwner.Uid {
			log.Fatalf("You do not own the file %s", dataFile)
		}
	}
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

	return readEntriesFromReader(f)
}

func readEntriesFromReader(r io.Reader) ([]PathEntry, error) {
	var entries []PathEntry
	scanner := bufio.NewScanner(r)
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
