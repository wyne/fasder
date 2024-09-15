<div align="center">

# Fasder - zoxide for files

<!--
[![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/wyne/fasder/total)](https://github.com/wyne/fasder/releases)
-->

Fasder let's you access **files and directories lightning quick**, inspired fasd.

It remembers which files and directories you use most frequently, so you can access
them in just a few keystrokes.<br />

[Installation](#installation) •
[Usage](#usage) •
[Why Fasder?](#why-fasder) •
[Features](#features)

![Demo](./demo.gif)

</div>

<hr />

Fasder is a modern reimagining of [clvv/fasd](http://github.com/clvv/fasd) that offers [zoxide](https://github.com/ajeetdsouza/zoxide)-style “frecent” (frequent + recent) access to files and directories.

Pronounced like “faster” with a ‘d’, Fasder tracks your most-used files and directories and lets you access them with minimal keystrokes. Need to reopen your `.zshrc`? Just type `v zsh` and you’re there.

### Key Benefits

- Fast access: Open your frequently-used files and directories with just a few characters.
- Minimal setup: Works out of the box with default aliases or customize it as you like.
- Powerful shortcuts: Built-in commands let you launch, edit, and navigate effortlessly.

### Examples

```bash
v def conf      # => vim /some/awkward/path/to/type/default.conf
j abc           # => cd /hell/of/a/awkward/path/to/get/to/abcdef
m movie         # => mplayer /path/to/awesome_movie.mp4
vim `f rc lo`   # => vim /etc/rc.local
```

## Installation

#### Basic Install:

```bash
brew install wyne/tap/fasder
echo 'eval "$(fasder --init auto)"' >> ~/.zshrc
```

#### Full Install (with dependencies and default aliases):

```bash
brew install wyne/tap/fasder fzf
echo 'eval "$(fasder --init auto aliases)"' >> ~/.zshrc
```

#### Migrate from `fasd`:

```bash
cp ~/.fasd ~/.fasder
```

## Usage

Pass `aliases` to init, example: `--init auto aliases`, to install these aliases:

```bash
a               # list files and directories
d               # directories only
f               # files only
v, vv           # open file in $EDITOR, vv for interactive
j, jj           # cd, jj for interactive
```

#### Example Commands

```bash
v def conf      # vim /awkward/path/default.conf
j abc           # cd /awkward/path/abcdef
vv foo          # Interactive file selection with fzf
jj foo          # Interactive directory navigation with fzf
```

The provided `v` and `vv` commands execute with program set in `$EDITOR`.
Configure with: `export EDITOR=nvim`.

## Base commands

```bash
fasder {query}        # files and directories
fasder -d {query}     # directories only
fasder -f {query}     # files only

{query} can be left empty to return all results
```

Example composition

```bash
alias a='fasder'        # both files and directories
alias d='fasder -d'     # directories only
alias f='fasder -f'     # files only
alias v='f -e $EDITOR'  # open file with $EDITOR
vim `f rc lo`           # on-the-fly command
```

See [shell.go](https://github.com/wyne/fasder/blob/main/shell.go) for provided aliases.

#### Options

```
fasder [options] [query ...]
  options:
        --init          Initialize fasder. Args: auto aliases
    -d, --directories   Dirs only
    -e, --exec {cmd}    Execute provided command against best match
    -f, --files         Files only
    -h, --help          Show this message
    -l, --list          List only. Omit rankings
    -R, --reverse       Reverse sort. Useful to pipe into fzf
    -s, --s             Show rank scores
    -v, --version       View version
```

## Matching

Matching works similarly to zoxide and obeys the following rules:

- The last word in the query must match the last segment of a path (split by "/" or ".").
  - `conf` will match `workspace/conf` but not `conf/project`
  - `conf yml` will match `config.yml` or `config/init.yml`
- Query words are matched in order to paths
  - `conf tmu` will match `config/tmux` but not `tmux/config.yml`.
- Path segment matches do not have to be adjacent
  - `work sub` will match `workspace/project/sub`

## Why fasder?

#### vs. zoxide

`zoxide` is great for directories. Fasder goes further—giving you quick access to both directories and files.

#### vs. fasd

`fasd` inspired `fasder`, but it’s now archived and written as a single, dense shell script. Fasder is built in a modern language, making it easier to read, maintain, and expand for more use cases.

# Features

- [x] Brew Formulae (`brew install wyne/tap/fasder`)
- [x] Aliases
- [ ] man page
- [ ] Shell Support
  - [x] Detect subshells
  - [x] zsh
    - [ ] autocomplete
  - [ ] bash
  - [ ] tcsh
- [ ] Ranking
  - [x] Shell hook to rank during normal operations
  - [x] Increment score on execution with -e flag
  - [x] Decay
- [ ] Matching
  - [x] Last segment matching
  - [x] Multiple path segment matching. Ex: {dir substr} {file substr} ([ref](https://github.com/clvv/fasd?tab=readme-ov-file#matching))
  - [ ] Full path matching. Ex: /some/dir/file
- [ ] Backends
  - [x] `fasd` format in `~/.fasder`
  - [x] neovim
    - [x] plugin [fasder.nvim](https://github.com/wyne/fasder.nvim)
    - [ ] shada
  - [ ] vim - viminfo
  - [ ] spotlight
  - [ ] recently used
- [ ] Flags
  - [x] `-R` reverse
  - [x] `-l` list paths without ranks
  - [x] `-f` files
  - [x] `-e` execute
  - [ ] `-t` recent access only
  - [ ] `-[0-9]` nth entry
  - [ ] `-b` only use backend
  - [ ] `-B` additional backend
  - [ ] `-i` interactive
