> [!WARNING]
> This project is under active development and may have breaking changes without warning.

# fasder - zoxide for files

This is a modern rewrite of [clvv/fasd](http://github.com/clvv/fasd) that provides
[zoxide](https://github.com/ajeetdsouza/zoxide)-style "frecently" used access to files and directories.

Fasder, pronounced like "faster" but with a d, offers quick access to commonly
used files and directories. Fasder tracks the files and directories that you access
and ranks them based on usage. You can then use the built in commands or
construct your own to reference them with minimal keystrokes.

For example, once you've opened your zsh config once, you can then use something
like `v zsh` or `v z` to immediately run `nvim ~/.zshrc`.

![Demo](./demo.gif)

## Installation

Bare installation (create your own aliases):

```bash
brew install wyne/tap/fasder
echo 'eval "$(fasder --init auto)"' >> ~/.zshrc
```

With dependencies and default aliases `f`, `a`, `d`, `v`, `vv`, `j`, `jj` (see [aliases](#aliases)):

```bash
brew install wyne/tap/fasder fzf
echo 'eval "$(fasder --init auto aliases)"' >> ~/.zshrc
```

Migrate from `fasd`:

```bash
cp ~/.fasd ~/.fasder
```

## Getting started

```bash
v def conf      # =>    vim /some/awkward/path/to/type/default.conf
j abc           # =>    cd /hell/of/a/awkward/path/to/get/to/abcdef
m movie         # =>    mplayer /whatever/whatever/whatever/awesome_movie.mp4
o eng paper     # =>    xdg-open /you/dont/remember/where/english_paper.pdf
vim `f rc lo`   # =>    vim /etc/rc.local
vim `f rc conf` # =>    vim /etc/rc.conf

v zsh           # =>    vim /commonly/used/file/.zshrc
vv foo          # =>    (interactive)
j foo           # =>    cd /commonly/used/path/foo
jj foo          # =>    (interactive)
j               # =>    cd - (back to previous directory)
```

The provided `v` and `vv` commands execute with program set in `$EDITOR`.
Configure with: `export EDITOR=nvim`.

## Usage

### Aliases

These aliases are installed by passing `aliases` as an init parameter.
Example: `--init auto aliases`.

```bash
alias a='fasder'        # both files and directories
alias d='fasder -d'     # directories only
alias f='fasder -f'     # files only

v zsh                   # =>    vim /commonly/used/file/.zshrc
vv zsh                  # =>    (interactive)
j foo                   # =>    cd /commonly/used/path/foo
jj foo                  # =>    (interactive)
j                       # =>    cd - (back to previous directory)
```

## Base commands

These commands will query the database or show the full database when no
query is provided. Results are ranked by usage.

```bash
fasder        # both files and directories
fasder -d     # directories only
fasder -f     # files only
```

Flags

```bash
fasder [options] [query ...]
  options:
        --init          Initialize fasder. Args: auto aliases
    -d, --directories   Dirs only
    -e, --exec <cmd>    Execute provided command against best match
    -f, --files         Files only
    -h, --help          Show this message
    -l, --list          List only. Omit rankings
    -R, --reverse       Reverse sort. Useful to pipe into fzf
    -s, --s             Show rank scores
    -v, --version       View version

fasder [--add|--sanitize,--proc] [paths ...]
  internal use:
        --sanitize      Internal: Sanitize command before processing
        --proc          Internal: Process a zsh-hook command
    -A, --add string    Internal: Add path to the store
```

## Compared to `zoxide`

[zoxide](https://github.com/ajeetdsouza/zoxide) only works on directories.
`Fasder` also works on files, allowing you to quickly access commonly used files as well.

## Compared to `fasd`

[clvv/fasd](http://github.com/clvv/fasd), the inspiration for `fasder`, has been
archived and will no longer be expanded. Additionally, it was written as
one large shell script which is difficult to read and maintain and contains no tests.

`fasder` is written in a modern language that is easy to adapt and expand to meet
more use cases.

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
  - [x] `-r` reverse
  - [x] `-l` list paths without ranks
  - [x] `-f` files
  - [x] `-e` execute
  - [ ] `-t` recent access only
  - [ ] `-[0-9]` nth entry
  - [ ] `-b` only use backend
  - [ ] `-B` additional backend
  - [ ] `-i` interactive
