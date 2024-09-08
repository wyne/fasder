> [!WARNING]
> This project is under very active development and likely to have breaking changes without warning.

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

# To-do

- [x] Support more aliases
- [x] VHS Tapes
- [x] Brew Formulae (`brew install wyne/tap/fasder`)
- [ ] Better ranking
  - [x] Increment score on execution with -e flag
  - [ ] Ranking decay
  - [ ] Remove entries from db on filtering
