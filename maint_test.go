package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/term"
)

func TestSubshellDetection(t *testing.T) {
	// Backup original stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Simulate non-terminal stdout
	if term.IsTerminal(int(w.Fd())) {
		t.Fatal("Expected non-terminal stdout")
	}

	// Create temporary files
	tempFile1, err := os.CreateTemp("", "file1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile1.Name())

	tempFile2, err := os.CreateTemp("", "file2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile2.Name())

	// Mock data for testing with real file paths
	mockData := fmt.Sprintf("%s|1.0|1627849200\n%s|2.0|1627849201", tempFile1.Name(), tempFile2.Name())

	// Use a temporary file for testing
	tempFile, err := os.CreateTemp("", "fasder_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write mock data to the temporary file
	if _, err := tempFile.WriteString(mockData); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// Set the environment variable to use the temporary file
	os.Setenv("_FASDER_DATA", tempFile.Name())
	defer os.Unsetenv("_FASDER_DATA")

	// Load the mock data file
	LoadFileStore()

	// Set up flags
	os.Args = []string{"cmd"}
	main()

	// Capture output
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Restore original stdout
	os.Stdout = oldStdout

	fmt.Printf("output %s", output)

	// Check if only one line is printed
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Errorf("Expected one line of output, but got %d lines", len(lines))
	}
}
