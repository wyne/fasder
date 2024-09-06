package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var dataFile string

// Struct to hold the file metadata
type FileEntry struct {
	Path         string
	Frequency    int
	LastAccessed time.Time
}

// Sorting Methods

type ByFrequencyThenRecency []FileEntry

func (a ByFrequencyThenRecency) Len() int      { return len(a) }
func (a ByFrequencyThenRecency) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByFrequencyThenRecency) Less(i, j int) bool {
	// Sort by frequency (descending), then by last accessed (descending)
	if a[i].Frequency != a[j].Frequency {
		return a[i].Frequency > a[j].Frequency
	}
	return a[i].LastAccessed.After(a[j].LastAccessed)
}

func sortEntries(entries []FileEntry) []FileEntry {
	sort.Sort(ByFrequencyThenRecency(entries))
	return entries
}

func openTopChoice(command string) {
	entries, err := readFileEntries()
	if err != nil {
		log.Fatal(err)
	}

	entries = sortEntries(entries) // Sort by frequency or use another sorting function

	if len(entries) > 0 {
		topEntry := entries[0]
		// Execute the specified command on the top entry
		cmd := exec.Command(command, topEntry.Path)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No entries found.")
	}
}

// Display

func displaySortedEntries(entries []FileEntry) {
	entries = sortEntries(entries)
	for _, entry := range entries {
		fmt.Printf("Path: %s, Frequency: %d, Last Accessed: %s\n",
			entry.Path, entry.Frequency, entry.LastAccessed.Format(time.RFC3339))
	}
}

// File methods

func LoadDataFile() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Silently return
		return
	}

	// Expand the ~ to the home directory
	dataFile = filepath.Join(homeDir, ".fasder")
}

// Reads the `.fasd` file and loads file entries into a slice
func readFileEntries() ([]FileEntry, error) {
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
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue // Skip malformed lines
		}

		freq, _ := strconv.Atoi(parts[1])
		lastAccessed, _ := time.Parse(time.RFC3339, parts[2])

		entry := FileEntry{
			Path:         parts[0],
			Frequency:    freq,
			LastAccessed: lastAccessed,
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

func logFileAccess(path string) {
	entries, err := readFileEntries()
	if err != nil {
		log.Fatal(err)
	}

	found := false
	for i, entry := range entries {
		if entry.Path == path {
			entries[i].Frequency++
			entries[i].LastAccessed = time.Now()
			found = true
			break
		}
	}

	if !found {
		// Add a new entry if the file hasn't been logged before
		newEntry := FileEntry{
			Path:         path,
			Frequency:    1,
			LastAccessed: time.Now(),
		}
		entries = append(entries, newEntry)
	}

	// Write updated entries back to the file
	writeToFile(entries)
}

func writeToFile(entries []FileEntry) {
	f, err := os.Create(dataFile) // Truncate and rewrite the file
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, entry := range entries {
		line := fmt.Sprintf("%s %d %s\n", entry.Path, entry.Frequency, entry.LastAccessed.Format(time.RFC3339))
		if _, err := f.WriteString(line); err != nil {
			log.Fatal(err)
		}
	}
}
