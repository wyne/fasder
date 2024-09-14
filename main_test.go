package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	flag "github.com/spf13/pflag"

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
