package logger

import (
	"log"
	"os"
	"path/filepath"
)

// Global variable to hold the logger
var Logger *log.Logger

// Initialize the logger
func InitLog() {
	// Open the log file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Silently return
		return
	}

	logFile := filepath.Join(homeDir, ".fasder.log")

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a new logger that writes to the file
	Logger = log.New(file, "", log.LstdFlags)

	// Optionally: log to both file and stdout
	// Logger = log.New(io.MultiWriter(file, os.Stdout), "", log.LstdFlags)
}
