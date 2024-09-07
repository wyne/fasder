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
	init := flag.Bool("init", false, "Initializers: zsh-hook")
	proc := flag.Bool("proc", false, "Process a command")
	display := flag.Bool("display", false, "Display sorted file entries")
	execCmd := flag.String("exec", "", "Command to open the top choice")

	LoadFileStore()

	flag.Parse()

	// Commands

	if *init {
		Init(flag.Args())
		return
	}

	if *proc {
		Proc(flag.Args())
		return
	}

	if *add != "" {
		Add(*add)
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
	entries, err := readFileStore()
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
