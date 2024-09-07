# fasder

This is a rewrite of [clvv/fasd](http://github.com/clvv/fasd) in go.

![Demo](./demo.gif)

## Installation

```bash
brew install wyne/tap/fasder
echo 'eval "$(fasder --init auto)"' >> ~/.zshrc
```

Built in aliases using `auto`:

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
- [ ] Increment score on execution with -e flag
- [ ] Better ranking
  - [ ] Ranking decay
- [ ] Remove from store on filtering
- [ ] VHS Tapes
- [ ] Brew Formulae
