package logger

import (
	"log"
	"os"
	"path/filepath"
)

// Global variable to hold the logger
var Log *log.Logger

type NoOpWriter struct{}

func (w *NoOpWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// Initialize the logger
func InitLog() {
	if len(os.Getenv("DEBUG")) > 0 {
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
		Log = log.New(file, "", log.LstdFlags)

		// Optionally: log to both file and stdout
		// Logger = log.New(io.MultiWriter(file, os.Stdout), "", log.LstdFlags)
	} else {
		Log = log.New(&NoOpWriter{}, "", 0)
	}
}
