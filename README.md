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
alias a='fasder'        # both files and directories
alias d='fasder -d'     # directories only
alias f='fasder -f'     # files only
alias z='fasder_cd -d'  # cd to best match. ex: `z work` to cd to workspace
```

## Usage

```bash
alias f='fasder -f' # files only
alias v='f -e nvim' # open in nvim
```

# To-do

- [x] Support more aliases
- [ ] Increment score on execution with -e flag
- [ ] Better ranking
- [ ] Remove from store on filtering
