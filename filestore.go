package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func readFileStore() ([]FileEntry, error) {
	var entries []FileEntry
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

		freq, _ := strconv.Atoi(parts[1])
		lastAccessed, _ := strconv.ParseInt(parts[2], 10, 64)

		entry := FileEntry{
			Path:         parts[0],
			Frequency:    freq,
			LastAccessed: lastAccessed,
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

func writeFileStore(entries []FileEntry) {
	f, err := os.Create(dataFile) // Truncate and rewrite the file
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, entry := range entries {
		line := fmt.Sprintf("%s|%d|%d\n", entry.Path, entry.Frequency, entry.LastAccessed)
		if _, err := f.WriteString(line); err != nil {
			log.Fatal(err)
		}
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
			entries[i].Frequency++
			entries[i].LastAccessed = time.Now().Unix()
			found = true
			break
		}
	}

	if !found {
		// Add a new entry if the file hasn't been logged before
		newEntry := FileEntry{
			Path:         path,
			Frequency:    1,
			LastAccessed: time.Now().Unix(),
		}
		entries = append(entries, newEntry)
	}

	// Write updated entries back to the file
	writeFileStore(entries)
}
