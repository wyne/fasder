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
	dataFile = os.Getenv("_FASDER_DATA")
	if dataFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Silently return
			return
		}
		// Expand the ~ to the home directory
		dataFile = filepath.Join(homeDir, ".fasder")
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

// writeFileStore ages the entries if the store has grown past the rank
// threshold, then writes them out. priorRank is the total read out of the data
// file: fasd sums the ranks of the lines it reads, before applying any
// rank + 1/rank increment and without counting the paths being added, so the
// caller supplies that rather than the post-update total.
func writeFileStore(entries []PathEntry, priorRank float64) {
	// Apply decay if cumulative rank exceeds threshold
	const threshold = 2000.0
	const decayFactor = 0.9

	if priorRank > threshold {
		logger.Log.Println("Rank threshold met. Decaying...")
		for i := range entries {
			entries[i].Rank *= decayFactor
		}
	}

	writeEntries(entries)
}

// writeEntries replaces the data file atomically, leaving every rank alone.
// Callers that must not rescore the store, such as deletion, use this directly:
// fasd's -D only filters the data file with sed and never ages what survives.
func writeEntries(entries []PathEntry) {
	mu.Lock() // Lock to prevent concurrent access
	defer mu.Unlock()

	tempPrefix := "fasder-"

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

// DeleteFromStore removes every entry whose path matches one of paths
func DeleteFromStore(paths []string) {
	if len(paths) == 0 {
		return
	}

	entries, err := readFileStore()
	if err != nil {
		log.Fatal(err)
	}

	toDelete := make(map[string]bool, len(paths))
	for _, path := range paths {
		toDelete[path] = true
	}

	remaining := make([]PathEntry, 0, len(entries))
	for _, entry := range entries {
		if !toDelete[entry.Path] {
			remaining = append(remaining, entry)
		}
	}

	if len(remaining) == len(entries) {
		// Nothing matched, so leave the file untouched
		return
	}

	writeEntries(remaining)
}

// pruneDecayed drops entries whose rank has fallen below 1. This mirrors the
// `$2 >= 1` pattern guarding fasd's add pass: lines under that threshold are
// never read into its table, so they are not written back out. Decay is what
// pushes entries under the line, and decayed values are written before being
// filtered, so an entry survives one pass past the decay that sank it, exactly
// as in fasd. Reads for queries are deliberately left unfiltered, matching
// fasd's query pass, which has no such guard.
func pruneDecayed(entries []PathEntry) []PathEntry {
	kept := make([]PathEntry, 0, len(entries))

	for _, entry := range entries {
		if entry.Rank >= 1 {
			kept = append(kept, entry)
		}
	}

	return kept
}

// AddToStore records the given paths in a single read-modify-write pass, the
// way fasd's --add handles its whole argument list in one awk invocation.
// Storing them one at a time would let the decay applied while writing the
// first path drop it below 1, so the prune at the head of the next path's pass
// would delete it before the operation finished. Repeated paths count once,
// matching fasd, whose BEGIN block keys new entries by path.
func AddToStore(paths ...string) {
	entries, err := readFileStore()
	if err != nil {
		log.Fatal(err)
	}

	entries = pruneDecayed(entries)

	// Totalled here, after the prune and before any increment, because that is
	// what fasd's `count += $2` accumulates: the ranks of the surviving lines
	// as they were read, with the paths being added excluded.
	var priorRank float64
	for _, entry := range entries {
		priorRank += entry.Rank
	}

	now := time.Now().Unix()
	seen := make(map[string]bool, len(paths))

	for _, path := range paths {
		if seen[path] {
			continue
		}
		seen[path] = true

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

				entries[i].LastAccessed = now
				found = true
				break
			}
		}

		if !found {
			// Add a new entry if the file hasn't been logged before
			entries = append(entries, PathEntry{
				Path:         path,
				Rank:         1,
				LastAccessed: now,
			})
		}
	}

	// Write updated entries back to the file
	writeFileStore(entries, priorRank)
}
