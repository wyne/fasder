package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestDisplayEntriesUsesFrecentScoreByDefault(t *testing.T) {
	oldStdout, r, w, err := redirectStdout()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdout = oldStdout
	}()

	displayEntries([]PathEntry{{
		Path:         "/recent-high-score",
		Rank:         10,
		LastAccessed: time.Now().Unix() - 10,
	}}, false, false)

	w.Close()
	output := captureOutput(r)
	if !strings.HasPrefix(output, "60.00000   /recent-high-score\n") {
		t.Fatalf("Expected frecency score output, got %q", output)
	}
}

func TestDisplayEntriesUsesRawRankScoreInRankMode(t *testing.T) {
	oldStdout, r, w, err := redirectStdout()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdout = oldStdout
	}()

	displayEntries([]PathEntry{{
		Path:         "/recent-high-score",
		Rank:         10,
		LastAccessed: time.Now().Unix() - 10,
	}}, false, true)

	w.Close()
	output := captureOutput(r)
	if !strings.HasPrefix(output, "10.00000   /recent-high-score\n") {
		t.Fatalf("Expected raw rank score output, got %q", output)
	}
}
