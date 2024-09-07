package main

import (
	"fmt"
)

// zshHook will setup a zshell hook to run and update based on commands
func zshHook() {
	fmt.Printf(`
# add zsh hook
_fasder_preexec() {
  { eval "fasder --proc $(fasder --sanitize $2)"; } >> /dev/null 2>&1
}
autoload -Uz add-zsh-hook
add-zsh-hook preexec _fasder_preexec
    `)
}
