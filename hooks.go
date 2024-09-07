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
	fmt.Print(`
alias a='fasder'
alias d='fasder -d'
alias f='fasder -f'
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
}
