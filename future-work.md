# Future work

This document tracks improvements that are useful but are not part of the
current implementation.

## Original fasd parity gaps

This list tracks the largest user-visible differences between Fasder and the
original shell script implementation in `~/workspace/clvv-fasd`.

### Highest impact

- [ ] Use true frecency scoring for normal queries.
      Original `fasd` scores normal results as rank multiplied by a recency
      weight. Fasder currently sorts by stored rank first and only uses
      last-access time as a tie-breaker.
- [ ] Add explicit query modes for rank-only and recent-only results.
      Original flags: `-r` for rank-only and `-t` for recent-only.
- [ ] Add `--delete` / `-D` to remove paths from the store.
- [ ] Add nth-result selection with `-[0-9]`.
      This supports commands such as selecting the second or third match
      without entering an interactive picker.
- [ ] Improve `-e` execution argument handling.
      Preserve command and path arguments without splitting quoted commands or
      paths incorrectly.

### Shell integration

- [ ] Support bash initialization.
      Original `fasd --init auto` detects bash and installs the shell hook plus
      command completion.
- [ ] Support POSIX shell initialization.
      Original `fasd` falls back to a prompt-based POSIX hook.
- [ ] Consider tcsh aliases and hooks if parity across older shell setups is
      still desired.
- [ ] Add command completion support.
      Original `fasd` supports command completion for bash and zsh.
- [ ] Add zsh word completion support.
      Original word completion expands markers such as `,query`, `f,query`,
      and `d,query` inside any command.

### Interactive workflows

- [ ] Add built-in interactive selection with `-i`.
      Fasder currently provides `jj` and `vv` fzf aliases, but original `fasd`
      can prompt interactively without relying on fzf.
- [ ] Revisit the default alias set.
      Original aliases include `s`, `sd`, `sf`, `z`, and `zz`; Fasder provides
      `a`, `d`, `f`, `v`, `j`, `jj`, and `vv`.
- [ ] Decide whether fzf-backed aliases should remain the preferred modern
      interaction model or coexist with original `-i` behavior.

### Backends

- [ ] Add backend selection flags.
      Original flags: `-b <name>` to use one backend and `-B <name>` to add a
      backend.
- [ ] Support the `current` backend for entries in the current directory.
- [ ] Support the `viminfo` backend or document a Neovim-first replacement.
- [ ] Support platform recent-file sources where practical.
      Original backends include Spotlight on macOS and GTK recently-used data
      on Linux.
- [ ] Consider custom backend hooks if user-defined sources are still useful.

### Configuration

- [ ] Add a read-only mode equivalent to `_FASD_RO`.
- [ ] Make current-directory tracking configurable.
      Original variable: `_FASD_TRACK_PWD`.
- [ ] Make the decay threshold configurable.
      Original variable: `_FASD_MAX`.
- [ ] Make fuzzy matching configurable.
      Original variable: `_FASD_FUZZY`.
- [ ] Decide whether Fasder should source rc files like `/etc/fasdrc` and
      `~/.fasdrc`, or keep configuration strictly environment-variable based.
- [ ] Add configurable logging/sink behavior if needed for debugging.

### Matching and compatibility checks

- [ ] Compare matching edge cases with original `fasd`.
      Include case-sensitive fallback, case-insensitive fallback, fuzzy
      fallback, trailing `/`, and trailing `$`.
- [ ] Add regression tests for subshell behavior and best-match output.
- [ ] Add tests that use migrated `~/.fasd` data copied to `~/.fasder`.
- [ ] Document intentional differences where Fasder should be modern rather
      than perfectly compatible.

## Command wrapper processing

Fasder's shell hook sends completed commands to `fasder --proc`. Before
processing a command, Fasder removes consecutive leading words listed in
`_FASDER_SHIFT`. The default list is currently:

```sh
sudo busybox
```

This works for simple commands such as `sudo vim file.txt`, but the
whitespace-separated list cannot describe wrapper syntax or consume wrapper
options.

### Expand the default wrapper list

- [ ] Add `doas`, a `sudo`-like privilege wrapper.
- [ ] Add `run0`, systemd's privilege-elevation command.
- [ ] Retain `sudo` and `busybox` for compatibility.

Proposed default after confirming platform behavior:

```sh
sudo doas run0 busybox
```

### Parse wrapper options

- [ ] Replace repeated word shifting with a wrapper-aware parser.
- [ ] Consume options and their values for `sudo`, `doas`, and `run0`.
- [ ] Preserve the command and its arguments after removing the wrapper.
- [ ] Fall back safely when wrapper syntax is incomplete or unknown.

Examples that should resolve to `vim file.txt`:

```sh
sudo -u root vim file.txt
doas -u root vim file.txt
run0 --user=root vim file.txt
```

### Support structured command runners

These should not be added to the flat `_FASDER_SHIFT` list. They require
recognizing a particular subcommand, delimiter, assignments, or options.

- [ ] Support `env`, including options and `NAME=value` assignments.
- [ ] Support `uv run`, without treating other `uv` subcommands as wrappers.
- [ ] Support `mise exec`, without treating other `mise` subcommands as
      wrappers.
- [ ] Evaluate `command`, `nohup`, `time`, and `nice` after wrapper-option
      parsing exists.

Examples:

```sh
env EDITOR=vim vim file.txt
uv run python script.py
mise exec -- node script.js
```

Adding `_FASDER_SHIFT="uv run"` does not model the `uv run` sequence today.
It independently marks both `uv` and `run` as removable leading words.

### Avoid ambiguous defaults

- [ ] Keep commands with common name collisions out of the default list.
- [ ] In particular, do not add `please` by default because that name is also
      used by the Please build system.
- [ ] Document how users can extend `_FASDER_SHIFT` for wrappers specific to
      their environment.

### Tests

- [ ] Cover each supported wrapper with and without options.
- [ ] Cover multiple nested wrappers, such as `sudo env FOO=bar command`.
- [ ] Verify that similarly named non-wrapper subcommands are not shifted.
- [ ] Verify that malformed wrapper invocations do not add incorrect paths.
- [ ] Compare supported cases with the original `fasd` behavior where
      applicable.
