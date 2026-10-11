package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func shellTestBinary(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fasder")
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", path, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	return path
}

func shellTestDir(t *testing.T) string {
	t.Helper()
	path, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func runShellTest(t *testing.T, shell, binary, script, cwd string, extraEnv ...string) string {
	t.Helper()
	args := []string{"-c", script}
	if shell == "bash" {
		args = []string{"--noprofile", "--norc", "-c", script}
	}
	cmd := exec.Command(shell, args...)
	cmd.Dir = cwd
	cmd.Env = append(os.Environ(),
		"PATH="+filepath.Dir(binary)+string(os.PathListSeparator)+os.Getenv("PATH"),
		"_FASDER_DATA="+filepath.Join(cwd, "data"),
		"BASH_ENV=", "ENV=", "HISTFILE=/dev/null")
	cmd.Env = append(cmd.Env, extraEnv...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", shell, err, output)
	}
	return string(output)
}

func TestBashHookPreservesPromptHandlers(t *testing.T) {
	for _, tc := range []struct{ name, setup, verify string }{
		{"unset", "unset PROMPT_COMMAND", `[ "$PROMPT_COMMAND" = _fasder_prompt_func ]`},
		{"empty", "PROMPT_COMMAND=''", `[ "$PROMPT_COMMAND" = _fasder_prompt_func ]`},
		{"string", "PROMPT_COMMAND='first_prompt; second_prompt'", `[ "$PROMPT_COMMAND" = '_fasder_prompt_func;first_prompt; second_prompt' ]`},
		{"array", "PROMPT_COMMAND=(first_prompt second_prompt)", `[ "${#PROMPT_COMMAND[@]}" -eq 3 ] && [ "${PROMPT_COMMAND[0]}" = _fasder_prompt_func ] && [ "${PROMPT_COMMAND[1]}" = first_prompt ] && [ "${PROMPT_COMMAND[2]}" = second_prompt ]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("bash", "--noprofile", "--norc", "-c", tc.setup+"\n"+BashHook()+"\n"+BashHook()+"\n"+tc.verify)
			cmd.Env = append(os.Environ(), "BASH_ENV=")
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("%v\n%s", err, output)
			}
		})
	}
}

func TestBashTrackingEndToEnd(t *testing.T) {
	binary := shellTestBinary(t)
	for _, tc := range []struct {
		name, command string
		tracked       bool
	}{
		{"quoted file", `cat 'tracked file'`, true},
		{"command flags", `cat -n 'tracked file'`, true},
		{"ignored command", `ls 'tracked file'`, false},
		{"blacklisted command", `cat --help 'tracked file'`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cwd := shellTestDir(t)
			path := filepath.Join(cwd, "tracked file")
			if err := os.WriteFile(path, nil, 0600); err != nil {
				t.Fatal(err)
			}
			script := `set -o history
history -c
HISTTIMEFORMAT='%F %T '
eval "$(fasder --init auto aliases)"
history -s "$TEST_COMMAND"; false; _fasder_prompt_func; status=$?
[ "$status" -eq 1 ] || exit 10
fasder -l
`
			got := runShellTest(t, "bash", binary, script, cwd, "TEST_COMMAND="+tc.command, "SHELL=/bin/zsh")
			if tc.tracked {
				lines := strings.Split(strings.TrimSpace(got), "\n")
				if len(lines) != 2 || !strings.Contains(got, cwd+"\n") || !strings.Contains(got, path+"\n") {
					t.Fatalf("unexpected tracked paths: %q", got)
				}
			} else if got != "" {
				t.Fatalf("ignored command tracked paths: %q", got)
			}
		})
	}
}

func TestBashAliasesEndToEnd(t *testing.T) {
	binary := shellTestBinary(t)
	cwd := shellTestDir(t)
	dest := filepath.Join(cwd, "destination space")
	if err := os.Mkdir(dest, 0700); err != nil {
		t.Fatal(err)
	}
	script := `shopt -s expand_aliases
eval "$(fasder --init auto aliases)"
fasder -A 'destination space'
d -l
j destination
pwd
j
`
	got := runShellTest(t, "bash", binary, script, cwd)
	want := dest + "\n" + dest + "\n" + cwd + "\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBashInteractivePromptTracksCommands(t *testing.T) {
	binary := shellTestBinary(t)
	cwd := shellTestDir(t)
	path := filepath.Join(cwd, "tracked file")
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("bash", "--noprofile", "--norc", "-i")
	cmd.Dir = cwd
	cmd.Env = append(os.Environ(),
		"PATH="+filepath.Dir(binary)+string(os.PathListSeparator)+os.Getenv("PATH"),
		"_FASDER_DATA="+filepath.Join(cwd, "data"),
		"BASH_ENV=", "HISTFILE=/dev/null", "HISTCONTROL=", "HISTIGNORE=", "PS1=")
	cmd.Stdin = strings.NewReader("eval \"$(fasder --init auto aliases)\"\ncat 'tracked file'\nfasder -l\nexit\n")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("interactive bash: %v", err)
	}
	if !strings.Contains(string(output), path+"\n") {
		t.Fatalf("file not tracked at prompt: %q", output)
	}
}

func TestAutoHookDetectsEvaluatingShell(t *testing.T) {
	binary := shellTestBinary(t)
	runShellTest(t, "bash", binary, `eval "$(fasder --init auto)"; declare -F _fasder_prompt_func >/dev/null && ! declare -F _fasder_preexec >/dev/null`, t.TempDir(), "SHELL=/bin/zsh")
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("zsh not installed")
	}
	runShellTest(t, "zsh", binary, `eval "$(fasder --init auto)"; (( $+functions[_fasder_preexec] )) && (( ! $+functions[_fasder_prompt_func] ))`, t.TempDir(), "SHELL=/bin/bash")
}
