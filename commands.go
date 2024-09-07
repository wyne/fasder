package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/wyne/fasder/logger"
)

// Initialization
func Init(args []string) {
	for _, initializer := range args {
		switch initializer {
		case "zsh-hook":
			zshHook()
		default:
			// fmt.Println("Unknown: ", initializer)
		}
	}
}

// Process command from shell hooks
func Proc(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Log.Println("Error getting working directory:", err)
		return
	}

	logger.Log.Printf("Processing: %s, %s", cwd, args)
}

// Sanitize command from shell hooks before processing
func Sanitize(args []string) {
	// Concatenate all arguments into a single string
	input := strings.Join(args, " ")

	// First, handle the command substitution: `$(...)` becomes `...`
	// This regex matches the command substitution and replaces it.
	reCommandSubstitution := regexp.MustCompile(`([^\\])\$\([^\)]*\)`)
	input = reCommandSubstitution.ReplaceAllString(input, "$1")

	// Then, replace special characters with a space: `|&;<>$`{}`
	reSpecialChars := regexp.MustCompile(`([^\\])[|&;<>$` + "`" + `{}]+`)
	input = reSpecialChars.ReplaceAllString(input, "$1 ")

	fmt.Printf("%s", input)
}

// Add an entry to the store
func Add(path string) {
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
