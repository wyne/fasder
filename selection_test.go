package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	flag "github.com/cornfeedhobo/pflag"
)

func TestNumericSelectionParsing(t *testing.T) {
	for _, tc := range []struct {
		args     []string
		want     uint
		wantExec string
		wantArgs []string
	}{
		{[]string{"-d12", "query"}, 12, "", []string{"query"}},
		{[]string{"-d012", "query"}, 12, "", []string{"query"}},
		{[]string{"-2", "-1"}, 1, "", []string{}},
		{[]string{"-1", "-0"}, 0, "", []string{}},
		{[]string{"--select=10"}, 10, "", []string{}},
		{[]string{"-e", "-2", "-3"}, 3, "-2", []string{}},
		{[]string{"--exec", "-2"}, 0, "-2", []string{}},
		{[]string{"-deprintf 12", "-2"}, 2, "printf 12", []string{}},
		{[]string{"-e=-2"}, 0, "-2", []string{}},
		{[]string{"--", "-2"}, 0, "", []string{"-2"}},
	} {
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		selection := registerSelectionFlag(flags)
		execCmd := flags.StringP("exec", "e", "", "")
		flags.BoolP("directories", "d", false, "")
		if err := flags.Parse(normalizeSelectionArgs(flags, tc.args)); err != nil {
			t.Fatal(err)
		}
		if *selection != tc.want || *execCmd != tc.wantExec || !reflect.DeepEqual(flags.Args(), tc.wantArgs) {
			t.Fatalf("%v: got selection=%d, exec=%q, args=%v", tc.args, *selection, *execCmd, flags.Args())
		}
	}
}

func TestNumericSelectionPreservesInternalCommands(t *testing.T) {
	for _, internal := range []string{"--proc", "--sanitize"} {
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Bool(internal[2:], false, "")
		selection := registerSelectionFlag(flags)
		args := []string{internal, "vim", "-2", "-d12"}
		if err := flags.Parse(normalizeSelectionArgs(flags, preserveCommandArgs(args))); err != nil {
			t.Fatal(err)
		}
		if *selection != 0 || !reflect.DeepEqual(flags.Args(), args[1:]) {
			t.Fatalf("command arguments changed: %v", flags.Args())
		}
	}
}

func TestNumericSelectionRejectsInvalidValues(t *testing.T) {
	for _, args := range [][]string{{"--select=-1"}, {"--select=abc"}, {"-999999999999999999999999999999"}} {
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		registerSelectionFlag(flags)
		if err := flags.Parse(normalizeSelectionArgs(flags, args)); err == nil {
			t.Fatalf("expected error for %v", args)
		}
	}
}

func TestNumericSelectionCLI(t *testing.T) {
	for _, tc := range []struct {
		name  string
		args  []string
		index int
	}{
		{"first", []string{"-1"}, 1},
		{"second", []string{"-2"}, 0},
		{"list override", []string{"-ls2"}, 0},
		{"reverse first", []string{"-R1"}, 0},
		{"reverse second", []string{"-R2"}, 1},
		{"execute", []string{"-2", "-e", "printf %s"}, 0},
		{"rank", []string{"-r1"}, 0},
		{"rank reverse", []string{"-rR1"}, 1},
		{"out of range", []string{"-3"}, -1},
		{"out of range execute", []string{"-3", "-e", "printf %s"}, -1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cleanup, r, w, paths := setupTest(t)
			defer cleanup()
			LoadFileStore()
			writeRankDisagreementStore(t, paths)
			os.Args = append([]string{"fasder"}, tc.args...)
			main()
			w.Close()
			output := captureOutput(r)
			want := ""
			if tc.index >= 0 {
				want = paths[tc.index]
				if tc.name != "execute" {
					want += "\n"
				}
			}
			if output != want {
				t.Fatalf("got %q, want %q", output, want)
			}
			entries, err := readFileStore()
			if err != nil {
				t.Fatal(err)
			}
			for i, entry := range entries {
				wantRank := []float64{10, 2}[i]
				if i == tc.index {
					wantRank += 1 / wantRank
				}
				if entry.Rank != wantRank {
					t.Fatalf("unexpected promotion: %v", entries)
				}
			}
		})
	}
}

func TestNumericSelectionMultiDigitAndFiltering(t *testing.T) {
	for _, tc := range []struct {
		args  []string
		index int
	}{
		{[]string{"-12"}, 1}, {[]string{"-d2"}, 10}, {[]string{"-f2"}, -1},
	} {
		cleanup, r, w, _ := setupTest(t)
		func() {
			defer cleanup()
			LoadFileStore()
			var entries []PathEntry
			for i := 0; i < 12; i++ {
				path := filepath.Join(t.TempDir(), fmt.Sprintf("match-%02d", i))
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
				entries = append(entries, PathEntry{Path: path, Rank: float64(i + 1), LastAccessed: 100})
			}
			file := filepath.Join(t.TempDir(), "file")
			if err := os.WriteFile(file, nil, 0600); err != nil {
				t.Fatal(err)
			}
			entries = append(entries, PathEntry{Path: file, Rank: 100, LastAccessed: 100})
			writeEntries(entries)
			os.Args = append([]string{"fasder"}, tc.args...)
			main()
			w.Close()
			want := ""
			if tc.index >= 0 {
				want = entries[tc.index].Path + "\n"
			}
			if got := captureOutput(r); got != want {
				t.Fatalf("%v: got %q, want %q", tc.args, got, want)
			}
		}()
	}
}
