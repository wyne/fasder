# fasder

This is a rewrite of [clvv/fasd](http://github.com/clvv/fasd) in go.

## Installation

```bash
brew install wyne/tap/fasder
```

Setup just the zsh-hook by adding to your `.zshrc`:

```bash
eval "$(fasder --init zsh-hook)"
```

You can optionally enable the builtin aliases:

```bash
eval "$(fasder --init zsh-hook aliases)"
```

```bash

alias a='fasder'    # both files and directories
alias d='fasder -d' # directories only
alias f='fasder -f' # files only

# function to execute built-in cd
fasder_cd() {
  if [ $# -le 1 ]; then
    fasder "$@"
  else
    local _fasder_ret="$(fasder -e 'printf %s' "$@")"
    [ -z "$_fasder_ret" ] && return
    [ -d "$_fasder_ret" ] && cd "$_fasder_ret" || printf %s\\n "$_fasder_ret"
  fi
}
alias z='fasder_cd -d'

```

## Usage

```bash
alias f='fasder -f' # files only
alias v='f -e nvim' # open in nvim
```

# To-do

- [ ] Increment score on execution with -e flag
- [ ] Support more aliases
- [ ] Better ranking
