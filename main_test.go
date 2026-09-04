package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	flag "github.com/cornfeedhobo/pflag"

	"golang.org/x/term"
)

func redirectStdout() (*os.File, *os.File, *os.File, error) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, nil, err
	}
	os.Stdout = w
	return oldStdout, r, w, nil
}

func setupTest(t *testing.T) (func(), *os.File, *os.File, []string) {
	originalArgs := os.Args
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	oldStdout, r, w, err := redirectStdout()
	if err != nil {
		t.Fatal(err)
	}

	if term.IsTerminal(int(w.Fd())) {
		t.Fatal("Expected non-terminal stdout")
	}

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
		os.Args = originalArgs
		os.Stdout = oldStdout
		os.Unsetenv("_FASDER_DATA")
		os.Remove(tempFile1.Name())
		os.Remove(tempFile2.Name())
		os.Remove(tempData.Name())
	}, r, w, []string{tempFile1.Name(), tempFile2.Name()}
}

func captureOutput(r *os.File) string {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func checkOutput(t *testing.T, r *os.File, expected []string) {
	output := captureOutput(r)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("Expected %v, but got %v", expected, lines)
	}
}

func TestList(t *testing.T) {
	teardown, r, w, paths := setupTest(t)
	defer teardown()

	LoadFileStore()
	os.Args = []string{"cmd", "-l"}
	main()

	w.Close()
	checkOutput(t, r, []string{paths[0], paths[1]})
}

func TestSubshellDetection(t *testing.T) {
	teardown, r, w, paths := setupTest(t)
	defer teardown()

	LoadFileStore()
	os.Args = []string{"cmd"}
	main()

	w.Close()
	checkOutput(t, r, []string{paths[1]})
}

func TestProcIgnoresDefaultCommands(t *testing.T) {
	_, readEntries := setupProcTest(t)

	Proc([]string{"ls", "tracked"})

	entries := readEntries()
	if len(entries) != 0 {
		t.Fatalf("Expected no entries, but got %v", entries)
	}
}

func TestProcBlacklistsAnyToken(t *testing.T) {
	_, readEntries := setupProcTest(t)

	if err := os.WriteFile("tracked", []byte("tracked"), 0600); err != nil {
		t.Fatal(err)
	}
	Proc([]string{"vim", "--help", "tracked"})

	entries := readEntries()
	if len(entries) != 0 {
		t.Fatalf("Expected no entries, but got %v", entries)
	}
}

func TestProcShiftsDefaultPrefixes(t *testing.T) {
	cwd, readEntries := setupProcTest(t)

	if err := os.WriteFile("tracked", []byte("tracked"), 0600); err != nil {
		t.Fatal(err)
	}
	expectedTracked, err := filepath.Abs("tracked")
	if err != nil {
		t.Fatal(err)
	}

	Proc([]string{"sudo", "vim", "tracked"})

	entries := readEntries()
	expectedPaths := []string{cwd, filepath.Clean(expectedTracked)}
	var actualPaths []string
	for _, entry := range entries {
		actualPaths = append(actualPaths, entry.Path)
	}
	if !reflect.DeepEqual(actualPaths, expectedPaths) {
		t.Fatalf("Expected %v, but got %v", expectedPaths, actualPaths)
	}
}

func TestAddFlagNormalizesMultiplePaths(t *testing.T) {
	originalArgs := os.Args
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	originalData := os.Getenv("_FASDER_DATA")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("one", []byte("one"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("nested", 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("two", []byte("two"), 0600); err != nil {
		t.Fatal(err)
	}

	dataPath := filepath.Join(tempDir, "fasder_data")
	if err := os.WriteFile(dataPath, nil, 0600); err != nil {
		t.Fatal(err)
	}
	os.Setenv("_FASDER_DATA", dataPath)
	defer func() {
		os.Args = originalArgs
		if originalData == "" {
			os.Unsetenv("_FASDER_DATA")
		} else {
			os.Setenv("_FASDER_DATA", originalData)
		}
		if err := os.Chdir(originalWd); err != nil {
			t.Fatal(err)
		}
	}()

	os.Args = []string{"cmd", "--add", "./one", "nested/../two"}
	main()

	entries, err := readFileStore()
	if err != nil {
		t.Fatal(err)
	}

	expectedOne, err := filepath.Abs("./one")
	if err != nil {
		t.Fatal(err)
	}
	expectedTwo, err := filepath.Abs("nested/../two")
	if err != nil {
		t.Fatal(err)
	}
	expectedPaths := []string{
		filepath.Clean(expectedOne),
		filepath.Clean(expectedTwo),
	}
	var actualPaths []string
	for _, entry := range entries {
		actualPaths = append(actualPaths, entry.Path)
	}
	if !reflect.DeepEqual(actualPaths, expectedPaths) {
		t.Fatalf("Expected %v, but got %v", expectedPaths, actualPaths)
	}
}

func setupProcTest(t *testing.T) (string, func() []PathEntry) {
	t.Helper()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dataPath := filepath.Join(tempDir, "fasder_data")
	if err := os.WriteFile(dataPath, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("_FASDER_DATA", dataPath)
	t.Setenv("_FASDER_BLACKLIST", "")
	t.Setenv("_FASDER_SHIFT", "")
	t.Setenv("_FASDER_IGNORE", "")
	LoadFileStore()

	t.Cleanup(func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Fatal(err)
		}
	})

	return cwd, func() []PathEntry {
		entries, err := readFileStore()
		if err != nil {
			t.Fatal(err)
		}
		return entries
	}
}
