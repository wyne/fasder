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
    alias v='fasder -f -e $EDITOR'
  `)

	// j - Jump to best match. If no arguments, jump to previous directory
	j := `
      j() {
        if [ "$#" -gt 0 ]; then
          cd "$(fasder -d -e 'printf %s' "$1")" || return 1
        else
          cd -
        fi
      }
    `
	fmt.Println(j)
}

func fzfAliases() {
	fmt.Println(`
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
  `)
	fmt.Println(`
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
  `)
}
