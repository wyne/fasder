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

func TestInternalCommandFlagsPreserveArguments(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantProc     bool
		wantSanitize bool
		wantVersion  bool
		wantExec     string
		wantArgs     []string
	}{
		{
			name:     "proc preserves a Fasder flag",
			args:     []string{"--proc", "echo", "--version"},
			wantProc: true,
			wantArgs: []string{"echo", "--version"},
		},
		{
			name:     "proc preserves an unknown flag",
			args:     []string{"--proc", "git", "status", "--short"},
			wantProc: true,
			wantArgs: []string{"git", "status", "--short"},
		},
		{
			name:     "proc equals form preserves an unknown flag",
			args:     []string{"--proc=true", "git", "status", "--short"},
			wantProc: true,
			wantArgs: []string{"git", "status", "--short"},
		},
		{
			name:     "proc preserves a value-consuming Fasder flag",
			args:     []string{"--proc", "vim", "-e", "foo"},
			wantProc: true,
			wantArgs: []string{"vim", "-e", "foo"},
		},
		{
			name:         "sanitize preserves command flags",
			args:         []string{"--sanitize", "vim", "--help"},
			wantSanitize: true,
			wantArgs:     []string{"vim", "--help"},
		},
		{
			name:     "existing delimiter remains unchanged",
			args:     []string{"--proc", "--", "echo", "--version"},
			wantProc: true,
			wantArgs: []string{"echo", "--version"},
		},
		{
			name:     "command delimiter survives injection",
			args:     []string{"--proc", "git", "log", "--", "README.md"},
			wantProc: true,
			wantArgs: []string{"git", "log", "--", "README.md"},
		},
		{
			name:     "proc without arguments",
			args:     []string{"--proc"},
			wantProc: true,
			wantArgs: []string{},
		},
		{
			name:     "delimiter before proc keeps it positional",
			args:     []string{"--", "--proc", "echo"},
			wantArgs: []string{"--proc", "echo"},
		},
		{
			name:        "normal Fasder flags are unchanged",
			args:        []string{"--version"},
			wantVersion: true,
			wantArgs:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := flag.NewFlagSet("fasder", flag.ContinueOnError)
			proc := flags.Bool("proc", false, "")
			sanitize := flags.Bool("sanitize", false, "")
			version := flags.Bool("version", false, "")
			execCmd := flags.StringP("exec", "e", "", "")

			if err := flags.Parse(preserveCommandArgs(tt.args)); err != nil {
				t.Fatal(err)
			}
			if *proc != tt.wantProc {
				t.Fatalf("Expected proc=%v, but got %v", tt.wantProc, *proc)
			}
			if *sanitize != tt.wantSanitize {
				t.Fatalf("Expected sanitize=%v, but got %v", tt.wantSanitize, *sanitize)
			}
			if *version != tt.wantVersion {
				t.Fatalf("Expected version=%v, but got %v", tt.wantVersion, *version)
			}
			if *execCmd != tt.wantExec {
				t.Fatalf("Expected exec=%q, but got %q", tt.wantExec, *execCmd)
			}
			if !reflect.DeepEqual(flags.Args(), tt.wantArgs) {
				t.Fatalf("Expected args %v, but got %v", tt.wantArgs, flags.Args())
			}
		})
	}
}

func TestProcDoesNotConsumeFasderFlags(t *testing.T) {
	originalArgs := os.Args
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	originalFlags := flag.CommandLine

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
	flag.CommandLine = flag.NewFlagSet("fasder", flag.ContinueOnError)
	flag.CommandLine.SetOutput(&bytes.Buffer{})
	os.Args = []string{"fasder", "--proc", "echo", "--version"}

	t.Cleanup(func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlags
		if err := os.Chdir(originalWd); err != nil {
			t.Fatal(err)
		}
	})

	main()

	entries, err := readFileStore()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Path != cwd {
		t.Fatalf("Expected proc to track %q, but got %v", cwd, entries)
	}
}
