package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

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

	tempDir1, err := os.MkdirTemp("", "dir1")
	if err != nil {
		t.Fatal(err)
	}

	mockData := fmt.Sprintf("%s|1.0|1627849200\n%s|2.0|1627849201\n%s|3.0|1627849202", tempFile1.Name(), tempFile2.Name(), tempDir1)

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
		os.Remove(tempDir1)
		os.Remove(tempData.Name())
	}, []string{tempFile1.Name(), tempFile2.Name(), tempDir1}
}

func TestPty(t *testing.T) {
	if os.Getenv("CI") != "" { // Check if running in a CI environment
		t.Skip("Skipping PTY test in CI environment")
	}

	teardown, paths := setupPtyTest(t)
	defer teardown()

	// Prepare the command to run the main program
	cmd := exec.Command("go", "run", ".", "-f") // This runs the current executable
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

func TestShell(t *testing.T) {
	if os.Getenv("CI") != "" { // Check if running in a CI environment
		t.Skip("Skipping PTY test in CI environment")
	}

	teardown, paths := setupPtyTest(t)
	defer teardown()

	pathSubStr := paths[2][len(paths[0])-3:]

	// Prepare the command to run the main program
	allCommands := []string{
		// ZshHook(),
		Aliases(),
		fmt.Sprintf("eval j var %v", pathSubStr),
		"eval pwd",
	}

	cmd := exec.Command("zsh", "-l", "-c", strings.Join(allCommands, "\n"))
	t.Logf("Command %s", strings.Join(cmd.Args, "\n"))

	// Create a pty
	pty, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("Failed to create pty: %v", err)
	}
	defer func() { _ = pty.Close() }() // Close the PTY after use

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
	expectedOutput := strings.Join([]string{paths[2]}, "\n")
	assert.Equal(t,
		strings.TrimSpace(expectedOutput),
		strings.TrimSpace(actualOutput),
		"Output mismatch: got '%v', want '%v'",
		actualOutput,
		expectedOutput)
}

// readFromPty reads the output from the pty until EOF or a timeout
func readFromPty(pty io.Reader) (string, error) {
	var buf bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Set a timeout
	defer cancel()

	done := make(chan error)

	go func() {
		_, err := io.Copy(&buf, pty) // Copy output from pty to buffer
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", err
		}
	case <-ctx.Done():
		return "", fmt.Errorf("read from pty timed out")
	}

	return buf.String(), nil // Return the captured output as a string
}
