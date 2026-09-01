package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wyne/fasder/logger"
)

const (
	defaultBlacklist = "--help"
	defaultShift     = "sudo busybox"
	defaultIgnore    = "fasder ls echo"
)

// Initialization
func Init(args []string) {
	for _, initializer := range args {
		switch initializer {
		case "auto":
			// TODO: Support other shells
			fmt.Println(ZshHook())
		case "zsh-hook":
			fmt.Println(ZshHook())
		case "aliases":
			fmt.Println(Aliases())
			fmt.Println(fzfAliases())
		}
	}
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

// Process command from shell hooks
func Proc(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Log.Println("Error getting working directory:", err)
		return
	}

	if containsAny(args, configuredWords("_FASDER_BLACKLIST", defaultBlacklist)) {
		return
	}

	shiftCommands := configuredWords("_FASDER_SHIFT", defaultShift)
	for len(args) > 0 && containsWord(shiftCommands, args[0]) {
		args = args[1:]
	}

	if len(args) > 0 && containsWord(configuredWords("_FASDER_IGNORE", defaultIgnore), args[0]) {
		return
	}

	paths := []string{cwd}
	if len(args) > 0 {
		paths = append(paths, args[1:]...)
	}
	Add(strings.Join(paths, " "))
}

func Add(args string) {
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
}

func configuredWords(name, defaultValue string) []string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
	}
	return strings.Fields(value)
}

func containsAny(values, candidates []string) bool {
	for _, value := range values {
		if containsWord(candidates, value) {
			return true
		}
	}
	return false
}

func containsWord(words []string, target string) bool {
	for _, word := range words {
		if word == target {
			return true
		}
	}
	return false
}
