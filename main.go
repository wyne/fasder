package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// Implement command-line flags
	search := flag.String("f", "", "Search for a file")
	add := flag.String("add", "", "Add object to store")
	display := flag.Bool("display", false, "Display sorted file entries")
	execCmd := flag.String("exec", "", "Command to open the top choice")

	flag.Parse()

	LoadStore()

	if *add != "" {
		logFileAccess(*add)
		return
	}

	if *search != "" {
		fmt.Println("Searching for file:", *search)
		// Add search logic here
		// print database
	} else {
		// print database
		fmt.Println("No search provided.")
	}

	// Retrieve and display sorted entries
	entries, err := readFromStore()
	if err != nil {
		log.Fatal(err)
	}

	if *display {
		// Sort entries based on the specified criteria
		displaySortedEntries(entries)
	}

	if *execCmd != "" {
		// Open the top choice with the specified command if -exec is set
		openTopChoice(*execCmd)
		// logFileAccess("/path/to/file1")
	}
}
