> [!WARNING]
> This project is under active development and may have breaking changes without warning.

# fasder

This is a rewrite of [clvv/fasd](http://github.com/clvv/fasd) in go.

![Demo](./demo.gif)

## Installation

```bash
brew install wyne/tap/fasder
echo 'eval "$(fasder --init auto)"' >> ~/.zshrc
```

Migrate from `fasd`:

```bash
cp .fasd .fasder
```

Built-in aliases using `auto`:

```bash
alias a='fasder'        # both files and directories
alias d='fasder -d'     # directories only
alias f='fasder -f'     # files only
alias z='fasder_cd -d'  # cd to best match. ex: `z work` to cd to workspace
alias j='fasder_cd -d'  # cd to best match. ex: `z work` to cd to workspace
```

FZF integration for quick selecting from a list

```bash
alias jj='fasder_cd -d'  # cd to best match. ex: `z work` to cd to workspace
```

## Custom Usage

```bash
alias v='f -e nvim'   # open best file match in nvim
```

# Advanced

```bash
# Search and select file with fzf, then execute with nvim
# Example: vv zsh
vv () {
  local file
  file="$(fasder -r -f -l "$1" | fzf -1 -0 --no-sort +m)"
  [ -n "$file" ] && echo "$file" | xargs nvim
}

# Search and select dir with fzf, then execute with nvim
# Example: jj work
jj () {
  local dir
  dir="$(fasder -r -d -l "$1" | fzf -1 -0 --no-sort +m)"  && cd "${dir}" || return 1
}
```

# Features

- [x] Brew Formulae (`brew install wyne/tap/fasder`)
- [ ] Shell Support
  - [x] zsh
  - [ ] bash
- [x] Aliases
- [ ] Ranking
  - [x] Shell hook to rank during normal operations
  - [x] Increment score on execution with -e flag
  - [ ] Decay
  - [ ] Remove entries from db on filtering
