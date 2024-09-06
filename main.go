package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// Implement command-line flags
	search := flag.String("f", "", "Search for a file")
	display := flag.Bool("display", false, "Display sorted file entries")
	execCmd := flag.String("exec", "", "Command to open the top choice")

	flag.Parse()

	LoadDataFile()

	if *search != "" {
		fmt.Println("Searching for file:", *search)
		// Add search logic here
	} else {
		fmt.Println("No search provided.")
	}

	// Retrieve and display sorted entries
	entries, err := readFileEntries()
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
		logFileAccess("/path/to/file1")
	}
}
