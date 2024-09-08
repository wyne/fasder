> [!WARNING]
> This project is under active development and may have breaking changes without warning.

# fasder

This is a rewrite of [clvv/fasd](http://github.com/clvv/fasd) in go.

Fasder, pronounced like "faster" but with a d, offers quick access to commonly
used files and directories. Fasder tracks the files and directories you access
and ranks them based on usage. You can then use the built in commands or
construct your own to reference them with minimal keystrokes.

For example, once you've opened your zsh config once, you can then use something
like `v .z` to immediately open `~/.zshrc` in neovim. See the aliases section
below to see how this works.

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

## Usage

### Built-in aliases and base commands

These aliases are installed with `auto` initializer, or individually
with `--init aliases`.

```bash
alias a='fasder'        # both files and directories
alias d='fasder -d'     # directories only
alias f='fasder -f'     # files only
```

Flags

- `-l` will omit scores and only print paths
- `-r` will reverse the list
- `-e {cmd}` will execute {cmd} on the best match

```bash
# Immediately open best match for query in $EDITOR
# Example: v {query}
# Leave query empty to open top ranked file from fasd -f
alias v='f -e $EDITOR'  # open best file match in $EDITOR
```

```bash
# Immediately cd to best match for query
# Example: j {query}
# Leave query empty to cd to top ranked dir from fasd -d
j() {
  cd "$(fasder -e 'printf %s' "$1")" || return 1
}
```

### Interactive Selection - requires [fzf](https://github.com/junegunn/fzf)

These aliases are installed with `auto` initializer or individually with
`--init fzf-aliases`.

```bash
# Interactive select from ranked files with fzf, then open in $EDITOR
# Example: vv {query}
# Leave query empty for full list
vv () {
  local selection
  # Get the selection from fasder and fzf
  selection=$(fasder -r -f -l "$1" | fzf -1 -0 --no-sort +m --height=10)

  # Check if a selection was made
  if [[ -n "$selection" ]]; then
      # Ensure the editor is set and handle potential issues
      if [[ -z "$EDITOR" ]]; then
          echo "EDITOR environment variable is not set."
          return 1
      fi

      # Use xargs with -r to prevent running the editor if no selection
      echo "Selection: $selection"
      echo "$selection" | xargs -r "$EDITOR"
  else
      echo "No selection made."
      return 1
  fi
}
```

```bash
# Interactive select from ranked files with fzf, then cd
# Example: jj {query}
# Leave query empty for full list
jj () {
  local selection
  selection=$(fasder -r -d -l "$1" | fzf -1 -0 --no-sort +m --height=10)
  if [[ -n "$selection" ]]; then
    echo "Selection: $selection"
    cd "$selection" || return 1
  else
    echo "No selection made"
    return 1
  fi
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
  - [x] Decay
  - [ ] Remove entries from file store on filtering
  - [ ] Full path search. Ex: {dir substr} {file substr}
