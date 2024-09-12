package main

import (
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
			name:       "Match only last path segment",
			searchTerm: "score",
			entries: []PathEntry{
				{Path: "/Users/justin/workspace/scorepad-react-native/android"},
				{Path: "/Users/justin/workspace/scorepad-react-native"},
			},
			filesOnly: false,
			dirsOnly:  false,
			expected: []PathEntry{
				{Path: "/Users/justin/workspace/scorepad-react-native"},
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
				t.Errorf("%s\n, fuzzyFind() = %v, want %v", tt.name, got, tt.expected)
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
