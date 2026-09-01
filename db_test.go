package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFuzzyFind(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		searchTerm string
		entries    []PathEntry
		filesOnly  bool
		dirsOnly   bool
		expected   []PathEntry
	}{
		{
			name:       "Gracefully handle empty paths",
			searchTerm: "one",
			entries: []PathEntry{
				{Path: "/one/two"},
				{Path: ""},
			},
			filesOnly: false,
			dirsOnly:  false,
			expected:  []PathEntry{},
		},
		{
			name:       "Match only last path segment",
			searchTerm: "one",
			entries: []PathEntry{
				{Path: "/one/two"},
				{Path: "/one"},
			},
			filesOnly: false,
			dirsOnly:  false,
			expected: []PathEntry{
				{Path: "/one"},
			},
		},
		{
			name:       "Match partial file name and extension",
			searchTerm: "tm con",
			entries: []PathEntry{
				{Path: "/Users/justin/.config/tmux/tmux.conf"},
				{Path: "/Users/justin/another/path.conf"},
			},
			filesOnly: false,
			dirsOnly:  false,
			expected: []PathEntry{
				{Path: "/Users/justin/.config/tmux/tmux.conf"},
			},
		},
		{
			name:       "Match partial term in path",
			searchTerm: ".conf tmu",
			entries: []PathEntry{
				{Path: "/Users/justin/.config/tmux/tmux.conf"},
				{Path: "/Users/justin/another/path.conf"},
			},
			filesOnly: false,
			dirsOnly:  false,
			expected: []PathEntry{
				{Path: "/Users/justin/.config/tmux/tmux.conf"},
			},
		},
		{
			name:       "Match with extra path segments",
			searchTerm: ".co azi.to",
			entries: []PathEntry{
				{Path: "/Users/justin/.config/yazi/yazi.toml"},
			},
			filesOnly: false,
			dirsOnly:  false,
			expected: []PathEntry{
				{Path: "/Users/justin/.config/yazi/yazi.toml"},
			},
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fuzzyFind(tt.entries, tt.searchTerm)
			if !equalPathEntries(got, tt.expected) {

				found := make([]string, len(got))
				for i, p := range got {
					found[i] = "    " + p.Path
				}

				want := make([]string, len(tt.expected))
				for i, p := range tt.expected {
					want[i] = "    " + p.Path
				}

				t.Errorf("%s\n\nSearch: %s\nFound:\n%v\nWant:\n%v",
					tt.name,
					tt.searchTerm,
					strings.Join(found, "\n"),
					strings.Join(want, "\n"))
			}
		})
	}
}

// Helper function to compare slices of PathEntry
func equalPathEntries(a, b []PathEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFrecentMultiplier(t *testing.T) {
	now := int64(1_700_000_000)
	tests := []struct {
		name         string
		lastAccessed int64
		expected     float64
	}{
		{
			name:         "within the last hour",
			lastAccessed: now - 3599,
			expected:     6,
		},
		{
			name:         "within the last day",
			lastAccessed: now - 86399,
			expected:     4,
		},
		{
			name:         "within the last week",
			lastAccessed: now - 604799,
			expected:     2,
		},
		{
			name:         "older than a week",
			lastAccessed: now - 604800,
			expected:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := frecentMultiplier(tt.lastAccessed, now)
			if got != tt.expected {
				t.Fatalf("Expected %v, but got %v", tt.expected, got)
			}
		})
	}
}

func TestSortEntriesUsesFrecentScore(t *testing.T) {
	now := int64(1_700_000_000)
	entries := []PathEntry{
		{Path: "/older-high-rank", Rank: 10, LastAccessed: now - 604800},
		{Path: "/recent-low-rank", Rank: 2, LastAccessed: now - 10},
		{Path: "/middle", Rank: 1, LastAccessed: now - 3600},
	}

	got := sortEntriesAt(entries, false, now)
	actualPaths := []string{got[0].Path, got[1].Path, got[2].Path}
	expectedPaths := []string{"/middle", "/older-high-rank", "/recent-low-rank"}
	if !reflect.DeepEqual(actualPaths, expectedPaths) {
		t.Fatalf("Expected %v, but got %v", expectedPaths, actualPaths)
	}
}

func TestSortEntriesReversesFrecentScore(t *testing.T) {
	now := int64(1_700_000_000)
	entries := []PathEntry{
		{Path: "/older-high-rank", Rank: 10, LastAccessed: now - 604800},
		{Path: "/recent-low-rank", Rank: 2, LastAccessed: now - 10},
		{Path: "/middle", Rank: 1, LastAccessed: now - 3600},
	}

	got := sortEntriesAt(entries, true, now)
	actualPaths := []string{got[0].Path, got[1].Path, got[2].Path}
	expectedPaths := []string{"/recent-low-rank", "/older-high-rank", "/middle"}
	if !reflect.DeepEqual(actualPaths, expectedPaths) {
		t.Fatalf("Expected %v, but got %v", expectedPaths, actualPaths)
	}
}
