package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const defaultBenchmarkEntries = 2000

type benchmarkTool struct {
	name string
	path string
	env  func(*benchmarkFixture) []string
}

type benchmarkFixture struct {
	cwd           string
	store         string
	seedStore     []byte
	addPath       string
	deletePath    string
	queryDirTerm  string
	queryFileTerm string
}

type benchmarkCase struct {
	name       string
	args       func(*benchmarkFixture) []string
	resetStore bool
}

func BenchmarkCliCompare(b *testing.B) {
	tools := benchmarkTools(b)
	entries := benchmarkEntryCount()
	cases := []benchmarkCase{
		{
			name: "query_broad_list",
			args: func(*benchmarkFixture) []string {
				return []string{"-l", "bench"}
			},
		},
		{
			name: "query_dir_filter",
			args: func(f *benchmarkFixture) []string {
				return []string{"-d", "-l", f.queryDirTerm}
			},
		},
		{
			name: "query_file_filter",
			args: func(f *benchmarkFixture) []string {
				return []string{"-f", "-l", f.queryFileTerm}
			},
		},
		{
			name: "subshell_best_match",
			args: func(*benchmarkFixture) []string {
				return []string{"bench"}
			},
			// Original fasd promotes the selected best match in this mode, so reset
			// the store before every iteration to keep each run comparable.
			resetStore: true,
		},
		{
			name: "add_new_path",
			args: func(f *benchmarkFixture) []string {
				return []string{"-A", f.addPath}
			},
			resetStore: true,
		},
		{
			name: "delete_path",
			args: func(f *benchmarkFixture) []string {
				return []string{"-D", f.deletePath}
			},
			resetStore: true,
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for _, tool := range tools {
				b.Run(tool.name, func(b *testing.B) {
					fixture := newBenchmarkFixture(b, entries)
					args := tc.args(fixture)
					env := tool.env(fixture)

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						if tc.resetStore {
							b.StopTimer()
							fixture.resetStore(b)
							b.StartTimer()
						}
						runBenchmarkCommand(b, tool.path, args, fixture.cwd, env)
					}
				})
			}
		})
	}
}

func benchmarkTools(b *testing.B) []benchmarkTool {
	b.Helper()

	var tools []benchmarkTool
	if path := os.Getenv("FASDER_BENCH_FASDER"); path != "" {
		tools = append(tools, benchmarkTool{
			name: "fasder",
			path: path,
			env: func(f *benchmarkFixture) []string {
				return benchmarkEnv(map[string]string{
					"_FASDER_DATA":  f.store,
					"_FASD_DATA":    "",
					"HOME":          filepath.Join(f.cwd, "..", "home"),
					"XDG_DATA_HOME": filepath.Join(f.cwd, "..", "xdg"),
					"PWD":           f.cwd,
				})
			},
		})
	}
	if path := os.Getenv("FASDER_BENCH_FASD"); path != "" {
		tools = append(tools, benchmarkTool{
			name: "fasd",
			path: path,
			env: func(f *benchmarkFixture) []string {
				return benchmarkEnv(map[string]string{
					"_FASD_AWK":      "awk",
					"_FASD_BACKENDS": "native",
					"_FASD_DATA":     f.store,
					"_FASD_FUZZY":    "2",
					"_FASD_MAX":      "2000",
					"_FASD_SINK":     "/dev/null",
					"_FASDER_DATA":   "",
					"HOME":           filepath.Join(f.cwd, "..", "home"),
					"PWD":            f.cwd,
				})
			},
		})
	}
	if len(tools) == 0 {
		b.Skip("set FASDER_BENCH_FASDER and/or FASDER_BENCH_FASD to benchmark CLI binaries")
	}
	for _, tool := range tools {
		if _, err := os.Stat(tool.path); err != nil {
			b.Fatalf("benchmark tool %s at %q is unavailable: %v", tool.name, tool.path, err)
		}
	}
	return tools
}

func benchmarkEntryCount() int {
	entries, err := strconv.Atoi(os.Getenv("FASDER_BENCH_ENTRIES"))
	if err != nil || entries < 2 {
		return defaultBenchmarkEntries
	}
	return entries
}

func newBenchmarkFixture(b *testing.B, entries int) *benchmarkFixture {
	b.Helper()

	root := b.TempDir()
	cwd := filepath.Join(root, "paths")
	if err := os.MkdirAll(cwd, 0700); err != nil {
		b.Fatal(err)
	}

	now := int64(1_700_000_000)
	lines := make([]string, 0, entries)
	var deletePath, queryDirTerm, queryFileTerm string
	for i := 0; i < entries; i++ {
		group := filepath.Join(cwd, fmt.Sprintf("project-%03d", i%100))
		if err := os.MkdirAll(group, 0700); err != nil {
			b.Fatal(err)
		}

		var path string
		if i%2 == 0 {
			queryDirTerm = fmt.Sprintf("bench-dir-%06d", i)
			path = filepath.Join(group, queryDirTerm)
			if err := os.MkdirAll(path, 0700); err != nil {
				b.Fatal(err)
			}
		} else {
			queryFileTerm = fmt.Sprintf("bench-file-%06d.txt", i)
			path = filepath.Join(group, queryFileTerm)
			if err := os.WriteFile(path, nil, 0600); err != nil {
				b.Fatal(err)
			}
		}
		if i == entries/2 {
			deletePath = path
		}

		rank := 1 + float64((i%100)+1)/10
		lastAccessed := now - int64(i*60)
		lines = append(lines, fmt.Sprintf("%s|%.5f|%d\n", path, rank, lastAccessed))
	}

	addPath := filepath.Join(cwd, "project-extra", "bench-dir-new")
	if err := os.MkdirAll(addPath, 0700); err != nil {
		b.Fatal(err)
	}

	fixture := &benchmarkFixture{
		cwd:           cwd,
		store:         filepath.Join(root, "store"),
		seedStore:     []byte(strings.Join(lines, "")),
		addPath:       addPath,
		deletePath:    deletePath,
		queryDirTerm:  queryDirTerm,
		queryFileTerm: queryFileTerm,
	}
	fixture.resetStore(b)
	return fixture
}

func (f *benchmarkFixture) resetStore(b *testing.B) {
	b.Helper()
	if err := os.WriteFile(f.store, f.seedStore, 0600); err != nil {
		b.Fatal(err)
	}
}

func benchmarkEnv(values map[string]string) []string {
	env := os.Environ()
	for key, value := range values {
		env = setEnv(env, key, value)
	}
	return env
}

func setEnv(env []string, key string, value string) []string {
	prefix := key + "="
	next := env[:0]
	for _, item := range env {
		if !strings.HasPrefix(item, prefix) {
			next = append(next, item)
		}
	}
	if value != "" {
		next = append(next, prefix+value)
	}
	return next
}

func runBenchmarkCommand(b *testing.B, path string, args []string, dir string, env []string) {
	b.Helper()

	cmd := exec.Command(path, args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		b.Fatalf("%s %s failed: %v", path, strings.Join(args, " "), err)
	}
}
