package main

import (
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	flag "github.com/cornfeedhobo/pflag"
)

func TestRecentScore(t *testing.T) {
	const now = int64(1700000000)
	for _, tc := range []struct {
		age  int64
		want float64
	}{
		{0, math.Sqrt(100000)}, {9, 100}, {99, math.Sqrt(1000)},
		{-1, math.Sqrt(100000)}, {-100, math.Sqrt(100000)},
	} {
		got := recentScore(PathEntry{Rank: 1000, LastAccessed: now - tc.age}, now)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("age %d: got %v, want %v", tc.age, got, tc.want)
		}
	}
}

func TestRecentSortIgnoresRankAndBreaksTiesByPath(t *testing.T) {
	const now = int64(1700000000)
	for _, reverse := range []bool{false, true} {
		entries := []PathEntry{
			{Path: "/z", Rank: 100, LastAccessed: now - 9},
			{Path: "/a", Rank: 1, LastAccessed: now - 9},
			{Path: "/old", Rank: 1000, LastAccessed: now - 99},
		}
		got := sortEntriesByRecentAt(entries, reverse, now)
		want := []string{"/old", "/a", "/z"}
		if reverse {
			want = []string{"/z", "/a", "/old"}
		}
		paths := []string{got[0].Path, got[1].Path, got[2].Path}
		if !reflect.DeepEqual(paths, want) {
			t.Fatalf("reverse %v: got %v, want %v", reverse, paths, want)
		}
	}
}

func TestSortFlagPrecedence(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want sortMode
	}{
		{[]string{"-rt"}, sortRecent}, {[]string{"-tr"}, sortRank},
		{[]string{"--recent", "--rank"}, sortRank}, {[]string{"-r", "-t"}, sortRecent},
		{[]string{"--recent=false"}, sortFrecent}, {[]string{"-t", "--recent=false"}, sortFrecent},
		{[]string{"--recent=invalid"}, sortFrecent},
	} {
		mode := sortFrecent
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		registerSortFlags(flags, &mode)
		err := flags.Parse(tc.args)
		if tc.args[0] == "--recent=invalid" {
			if err == nil {
				t.Fatal("expected invalid boolean error")
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if mode != tc.want {
			t.Fatalf("%v: got %v, want %v", tc.args, mode, tc.want)
		}
	}
}

func TestDisplayRecentScore(t *testing.T) {
	old, r, w, err := redirectStdout()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { os.Stdout = old; r.Close() }()
	displayEntriesAt([]PathEntry{{Path: "/recent", Rank: 1000, LastAccessed: 91}}, false, sortRecent, 100)
	w.Close()
	if got := captureOutput(r); got != "100.00000  /recent\n" {
		t.Fatalf("got %q", got)
	}
}

func TestRecentFlagCLI(t *testing.T) {
	for _, tc := range []struct {
		name     string
		args     []string
		selected bool
		reverse  bool
	}{
		{"list", []string{"-tl"}, false, false},
		{"reverse list", []string{"-tlR"}, false, true},
		{"subshell", []string{"-t"}, true, false},
		{"reverse subshell", []string{"-tR"}, true, false},
		{"execute", []string{"-te", "printf %s"}, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cleanup, r, w, paths := setupTest(t)
			defer cleanup()
			LoadFileStore()
			writeEntries([]PathEntry{
				{Path: paths[0], Rank: 100, LastAccessed: time.Now().Unix() - 99},
				{Path: paths[1], Rank: 2, LastAccessed: time.Now().Unix() - 9},
			})
			os.Args = append([]string{"fasder"}, tc.args...)
			main()
			w.Close()
			want := paths
			if tc.reverse {
				want = []string{paths[1], paths[0]}
			}
			if tc.selected {
				want = []string{paths[1]}
			}
			checkOutput(t, r, want)
			entries, err := readFileStore()
			if err != nil {
				t.Fatal(err)
			}
			wantRank := 2.0
			if tc.selected {
				wantRank = 2.5
			}
			if entries[0].Rank != 100 || entries[1].Rank != wantRank {
				t.Fatalf("unexpected promotion: %v", entries)
			}
		})
	}
}
