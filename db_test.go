package main

import (
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
