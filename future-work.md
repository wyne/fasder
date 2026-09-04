# Future work

This document tracks improvements that are useful but are not part of the
current implementation.

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
