package main

import (
	"fmt"
)

// zshHook will setup a zshell hook to run and update based on commands
func zshHook() {
	fmt.Printf(`
    _fasder_preexec() {
      { eval "fasder --proc $(fasder --sanitize $2)"; } >> /dev/null 2>&1
    }
    autoload -Uz add-zsh-hook
    add-zsh-hook preexec _fasder_preexec
  `)
}

func aliases() {
	fmt.Println(`
    alias a='fasder'
    alias d='fasder -d'
    alias f='fasder -f'
  `)

	// j - Jump to best match. If no arguments, jump to previous directory
	fmt.Println(`
    j() {
      cd "$(fasder -e 'printf %s' "$1")" || return 1
    }  
  `)
}

func fzfAliases() {
	fmt.Println(`
    jj() {
      cd "$(fasder -r -d -l "$1" | fzf -1 -0 --no-sort +m)" || return 1
    }
  `)
	fmt.Println(`
    vv() {
      fasder -r -f -l "$1" | fzf -1 -0 --no-sort +m | xargs -r $EDITOR
    }
  `)
}
