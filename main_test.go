package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

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

func TestSubshellDetectionReversed(t *testing.T) {
	teardown, r, w, paths := setupTest(t)
	defer teardown()

	LoadFileStore()
	os.Args = []string{"cmd", "-R"}
	main()

	w.Close()
	// -R reverses display order, but the single best match stays the same.
	checkOutput(t, r, []string{paths[1]})
}

func TestExecuteUsesBestMatchWhenReversed(t *testing.T) {
	teardown, r, w, paths := setupTest(t)
	defer teardown()

	LoadFileStore()
	os.Args = []string{"cmd", "-R", "-e", "echo"}
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
	_, readEntries := setupProcTest(t)

	Proc([]string{"sudo", "ls", "tracked"})

	if entries := readEntries(); len(entries) != 0 {
		t.Fatalf("Expected no entries, but got %v", entries)
	}
}

func TestProcUsesConfiguredWords(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		value string
		args  []string
	}{
		{
			name:  "blacklist",
			env:   "_FASDER_BLACKLIST",
			value: "blocked",
			args:  []string{"vim", "blocked", "tracked"},
		},
		{
			name:  "shift",
			env:   "_FASDER_SHIFT",
			value: "doas",
			args:  []string{"doas", "ls", "tracked"},
		},
		{
			name:  "ignore",
			env:   "_FASDER_IGNORE",
			value: "vim",
			args:  []string{"vim", "tracked"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, readEntries := setupProcTest(t)
			t.Setenv(tt.env, tt.value)

			Proc(tt.args)

			if entries := readEntries(); len(entries) != 0 {
				t.Fatalf("Expected no entries, but got %v", entries)
			}
		})
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
	if err := os.WriteFile("tracked", []byte("tracked"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("_FASDER_DATA", dataPath)
	flag.CommandLine = flag.NewFlagSet("fasder", flag.ContinueOnError)
	flag.CommandLine.SetOutput(&bytes.Buffer{})
	os.Args = []string{"fasder", "--proc", "vim", "--version", "tracked"}

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
	expectedPaths := []string{cwd, filepath.Join(cwd, "tracked")}
	var actualPaths []string
	for _, entry := range entries {
		actualPaths = append(actualPaths, entry.Path)
	}
	if !reflect.DeepEqual(actualPaths, expectedPaths) {
		t.Fatalf("Expected proc to track %v, but got %v", expectedPaths, actualPaths)
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
	// Empty means use the built-in defaults, pinning these tests to known lists.
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

// Writes a script that records each argument it receives, one per line, to the
// file named by FASDER_TEST_ARGV. Returns the script path.
func writeArgvRecorder(t *testing.T, dir string, logPath string) string {
	t.Helper()

	script := filepath.Join(dir, "record-argv")
	body := "#!/bin/sh\n: > \"$FASDER_TEST_ARGV\"\nfor arg in \"$@\"; do printf '%s\\n' \"$arg\" >> \"$FASDER_TEST_ARGV\"; done\n"
	if err := os.WriteFile(script, []byte(body), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FASDER_TEST_ARGV", logPath)
	return script
}

func recordedArgs(t *testing.T, logPath string) []string {
	t.Helper()

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

func TestExecutePassesSpacedPathAsSingleArgument(t *testing.T) {
	cwd, readEntries := setupProcTest(t)

	spaced := filepath.Join(cwd, "my file.txt")
	if err := os.WriteFile(spaced, nil, 0600); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(cwd, "argv.log")
	script := writeArgvRecorder(t, cwd, logPath)

	execute([]PathEntry{{Path: spaced, Rank: 1, LastAccessed: time.Now().Unix()}}, script, false)

	expectedArgs := []string{spaced}
	if got := recordedArgs(t, logPath); !reflect.DeepEqual(got, expectedArgs) {
		t.Fatalf("Expected %v, but got %v", expectedArgs, got)
	}

	// The rank increment also has to survive the space.
	entries := readEntries()
	if len(entries) != 1 || entries[0].Path != spaced {
		t.Fatalf("Expected the spaced path to be tracked, but got %v", entries)
	}
}

func TestExecuteSplitsMultiWordCommand(t *testing.T) {
	cwd, _ := setupProcTest(t)

	spaced := filepath.Join(cwd, "my file.txt")
	if err := os.WriteFile(spaced, nil, 0600); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(cwd, "argv.log")
	script := writeArgvRecorder(t, cwd, logPath)

	// The j alias runs `fasder -de 'printf %s'`, so the command itself must
	// keep word-splitting even while the path stays whole.
	execute([]PathEntry{{Path: spaced, Rank: 1, LastAccessed: time.Now().Unix()}}, script+" --flag", false)

	expectedArgs := []string{"--flag", spaced}
	if got := recordedArgs(t, logPath); !reflect.DeepEqual(got, expectedArgs) {
		t.Fatalf("Expected %v, but got %v", expectedArgs, got)
	}
}

// storedPaths returns the paths currently in the store, in file order.
func storedPaths(t *testing.T, readEntries func() []PathEntry) []string {
	t.Helper()

	var paths []string
	for _, entry := range readEntries() {
		paths = append(paths, entry.Path)
	}
	return paths
}

func seedStore(t *testing.T, cwd string, names ...string) []string {
	t.Helper()

	var seeded []string
	for _, name := range names {
		path := filepath.Join(cwd, name)
		if err := os.MkdirAll(path, 0700); err != nil {
			t.Fatal(err)
		}
		AddToStore(path)
		seeded = append(seeded, path)
	}
	return seeded
}

func runMain(t *testing.T, args ...string) {
	t.Helper()

	originalArgs := os.Args
	originalFlags := flag.CommandLine
	t.Cleanup(func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlags
	})

	flag.CommandLine = flag.NewFlagSet("fasder", flag.ContinueOnError)
	flag.CommandLine.SetOutput(&bytes.Buffer{})
	os.Args = append([]string{"fasder"}, args...)
	main()
}

func TestDeleteAcceptsMultiplePositionalPaths(t *testing.T) {
	cwd, readEntries := setupProcTest(t)
	seeded := seedStore(t, cwd, "one", "two", "three")

	// fasd's usage is `fasd [-A|-D] [paths ...]`, so every trailing argument
	// is a path, not just the one attached to the flag.
	runMain(t, "-D", seeded[0], seeded[1])

	expected := []string{seeded[2]}
	if got := storedPaths(t, readEntries); !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestDeleteAcceptsRepeatedFlags(t *testing.T) {
	cwd, readEntries := setupProcTest(t)
	seeded := seedStore(t, cwd, "one", "two", "three")

	runMain(t, "-D", seeded[0], "-D", seeded[1])

	expected := []string{seeded[2]}
	if got := storedPaths(t, readEntries); !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestDeleteNormalizesPaths(t *testing.T) {
	cwd, readEntries := setupProcTest(t)
	seedStore(t, cwd, "one", "two", "three")

	// Relative, dot-prefixed, and trailing-slash forms all have to resolve to
	// the stored absolute path, the way fasd's normalization sed does.
	runMain(t, "-D", "./one", "two/", filepath.Join("one", "..", "three"))

	if got := storedPaths(t, readEntries); len(got) != 0 {
		t.Fatalf("Expected an empty store, but got %v", got)
	}
}

func TestDeleteCurrentDirectoryShorthand(t *testing.T) {
	cwd, readEntries := setupProcTest(t)
	seeded := seedStore(t, cwd, "one", "two")

	if err := os.Chdir(seeded[0]); err != nil {
		t.Fatal(err)
	}
	runMain(t, "-D", ".")

	expected := []string{seeded[1]}
	if got := storedPaths(t, readEntries); !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestDeleteRemovesStaleEntry(t *testing.T) {
	cwd, readEntries := setupProcTest(t)
	seeded := seedStore(t, cwd, "one", "two")

	// Clearing entries for paths that no longer exist is the main reason to
	// run this, so deletion must not check the filesystem.
	if err := os.Remove(seeded[0]); err != nil {
		t.Fatal(err)
	}
	runMain(t, "-D", seeded[0])

	expected := []string{seeded[1]}
	if got := storedPaths(t, readEntries); !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestDeleteUnknownPathLeavesStoreIntact(t *testing.T) {
	cwd, readEntries := setupProcTest(t)
	seeded := seedStore(t, cwd, "one", "two")

	runMain(t, "-D", filepath.Join(cwd, "never-tracked"))

	if got := storedPaths(t, readEntries); !reflect.DeepEqual(got, seeded) {
		t.Fatalf("Expected %v, but got %v", seeded, got)
	}
}

func TestDeleteDoesNotDecaySurvivingRanks(t *testing.T) {
	cwd, readEntries := setupProcTest(t)
	seeded := seedStore(t, cwd, "busy", "alsobusy", "stale")

	// Push the survivors past writeFileStore's 2000 decay threshold. fasd's -D
	// only filters the data file with sed, so deleting one path must not
	// rescore the records that remain.
	writeEntries([]PathEntry{
		{Path: seeded[0], Rank: 1100, LastAccessed: 100},
		{Path: seeded[1], Rank: 1000, LastAccessed: 200},
		{Path: seeded[2], Rank: 100, LastAccessed: 300},
	})

	runMain(t, "-D", seeded[2])

	entries := readEntries()
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, but got %v", entries)
	}
	if entries[0].Rank != 1100 || entries[1].Rank != 1000 {
		t.Fatalf("Expected ranks 1100 and 1000 to survive untouched, but got %v and %v",
			entries[0].Rank, entries[1].Rank)
	}
}

func TestAddToStoreDropsDecayedEntries(t *testing.T) {
	cwd, readEntries := setupProcTest(t)

	faded := filepath.Join(cwd, "faded")
	fresh := filepath.Join(cwd, "fresh")
	for _, dir := range []string{faded, fresh} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	writeEntries([]PathEntry{
		{Path: faded, Rank: 0.05, LastAccessed: time.Now().Unix()},
		{Path: filepath.Join(cwd, "boundary"), Rank: 1, LastAccessed: time.Now().Unix()},
	})

	AddToStore(fresh)

	var got []string
	for _, entry := range readEntries() {
		got = append(got, entry.Path)
	}
	expected := []string{filepath.Join(cwd, "boundary"), fresh}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %v, but got %v", expected, got)
	}
}

func TestAddToStoreResetsDecayedEntryInsteadOfBoosting(t *testing.T) {
	cwd, readEntries := setupProcTest(t)

	faded := filepath.Join(cwd, "faded")
	if err := os.MkdirAll(faded, 0700); err != nil {
		t.Fatal(err)
	}
	writeEntries([]PathEntry{{Path: faded, Rank: 0.05, LastAccessed: time.Now().Unix()}})

	// rank + 1/rank would send 0.05 to 20.05, leapfrogging genuinely frequent
	// paths. fasd never sees a sub-1 rank here, so re-adding starts over at 1.
	AddToStore(faded)

	entries := readEntries()
	if len(entries) != 1 || entries[0].Rank != 1 {
		t.Fatalf("Expected a single entry with rank 1, but got %v", entries)
	}
}

func TestQueriesStillSeeDecayedEntries(t *testing.T) {
	cwd, _ := setupProcTest(t)

	faded := filepath.Join(cwd, "faded")
	if err := os.MkdirAll(faded, 0700); err != nil {
		t.Fatal(err)
	}
	writeEntries([]PathEntry{{Path: faded, Rank: 0.05, LastAccessed: time.Now().Unix()}})

	// fasd's query pass has no `$2 >= 1` guard, so reads stay unfiltered and a
	// decayed entry remains visible until the next add sweeps it away.
	entries, err := readFileStore()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Path != faded {
		t.Fatalf("Expected the decayed entry to still be readable, but got %v", entries)
	}
}
