package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/creack/pty"
	"github.com/stretchr/testify/assert"
)

func setupPtyTest(t *testing.T) (func(), []string) {
	tempFile1, err := os.CreateTemp("", "file1")
	if err != nil {
		t.Fatal(err)
	}

	tempFile2, err := os.CreateTemp("", "file2")
	if err != nil {
		t.Fatal(err)
	}

	mockData := fmt.Sprintf("%s|1.0|1627849200\n%s|2.0|1627849201", tempFile1.Name(), tempFile2.Name())

	tempData, err := os.CreateTemp("", "fasder_test")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tempData.WriteString(mockData); err != nil {
		t.Fatal(err)
	}
	tempData.Close()

	os.Setenv("_FASDER_DATA", tempData.Name())

	return func() {
		os.Unsetenv("_FASDER_DATA")
		os.Remove(tempFile1.Name())
		os.Remove(tempFile2.Name())
		os.Remove(tempData.Name())
	}, []string{tempFile1.Name(), tempFile2.Name()}

}

func TestPty(t *testing.T) {
	if os.Getenv("CI") != "" { // Check if running in a CI environment
		t.Skip("Skipping PTY test in CI environment")
	}

	teardown, paths := setupPtyTest(t)
	defer teardown()

	// Prepare the command to run the main program
	cmd := exec.Command("go", "run", ".") // This runs the current executable
	t.Logf("Command %s", os.Args[0])

	// Create a pty
	pty, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("Failed to create pty: %v", err)
	}
	defer pty.Close() // Ensure pty is closed after the test

	// Read output from the pty
	output, err := readFromPty(pty)
	if err != nil {
		t.Fatalf("Failed to read from pty: %v", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		t.Fatalf("Command finished with error: %v", err)
	}

	// Normalize line endings for comparison
	actualOutput := strings.ReplaceAll(output, "\r\n", "\n")

	// Assert the output against the expectation
	expectedOutput := strings.Join([]string{paths[0], paths[1]}, "\n")
	assert.Equal(t,
		strings.TrimSpace(expectedOutput),
		strings.TrimSpace(actualOutput),
		"Output mismatch: got '%v', want '%v'",
		actualOutput,
		expectedOutput)

}

// readFromPty reads the output from the pty until EOF
func readFromPty(pty io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, pty) // Copy output from pty to buffer
	if err != nil {
		return "", err
	}
	return buf.String(), nil // Return the captured output as a string
}
