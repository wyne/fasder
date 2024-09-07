package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wyne/fasder/logger"
)

// Initialization
func Init(args []string) {
	for _, initializer := range args {
		switch initializer {
		case "auto":
			zshHook()
			aliases()
			fzfAliases()
		case "zsh-hook":
			zshHook()
		case "aliases":
			aliases()
		case "fzf-aliases":
			fzfAliases()
		}
	}
}

// Sanitize command from shell hooks before processing
func Sanitize(args []string) {
	// logger.Log.Printf("--sanitize: %s", strings.Join(args, " "))

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

// Process command from shell hooks
func Proc(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Log.Println("Error getting working directory:", err)
		return
	}

	// TODO: ignores
	// TODO: blacklists
	// TODO: shifts?

	// logger.Log.Printf("--proc: %s %s", cwd, strings.Join(args, " "))

	Add(fmt.Sprintf("%s %s", cwd, strings.Join(args, " ")))
}

func Add(args string) {
	// logger.Log.Printf("--add: %s", args)

	var validPaths []string

	// Iterate over the arguments and validate paths
	for _, arg := range strings.Split(args, " ") {
		if _, err := os.Stat(arg); err == nil {
			validPaths = append(validPaths, arg)
		}
	}

	// Convert paths to absolute form and simplify
	var absolutePaths []string
	for _, path := range validPaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("Error converting path to absolute form: %v\n", err)
			continue
		}
		// Simplify the path
		cleanPath := filepath.Clean(absPath)
		absolutePaths = append(absolutePaths, cleanPath)
	}

	for _, path := range absolutePaths {
		AddToStore(path)
	}

	// Join paths with a '|'
	result := strings.Join(absolutePaths, "|")

	logger.Log.Printf("--add paths: %s", result)
}
