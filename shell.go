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

	fmt.Println(`
    fz() {
      local dir
      dir="$(fasder -r -d -l "$1" | fzf -1 -0 --no-sort +m)" && cd "${dir}" || return 1
    }
  `)

	fmt.Println(`
    fasder_cd() {
      if [ $# -le 1 ]; then
        fasder -l "$@"
      else
        local _fasder_ret="$(fasder -e 'printf %s' "$@")"
        [ -z "$_fasder_ret" ] && return
        [ -d "$_fasder_ret" ] && cd "$_fasder_ret" || printf %s\n "$_fasder_ret"
      fi
    }
    alias z='fasder_cd -d'
  `)

	// # j - same as z, but if no arguments, jump to previous directory
	fmt.Println(`
    j() {
      if [ "$#" -gt 0 ]; then
        fasder_cd -d $1
      else
        cd -
      fi
    }
  `)
}

func fzfAliases() {
	fmt.Println(`
    jj() {
      local dir
      dir="$(fasder -r -d -l "$1" | fzf -1 -0 --no-sort +m)" && cd "${dir}" || return 1
    }
  `)
}
