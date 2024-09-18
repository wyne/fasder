package main

// ZshHook will setup a zshell hook to run and update based on commands
func ZshHook() string {
	return `
    _fasder_preexec() {
      { eval "fasder --proc $(fasder --sanitize $2)"; } >> /dev/null 2>&1
    }
    autoload -Uz add-zsh-hook
    add-zsh-hook preexec _fasder_preexec
    `
}

func Aliases() string {
	return `
    alias a='fasder'
    alias d='fasder -d'
    alias f='fasder -f'
    alias v='fasder -fe $EDITOR'
	j() {
	    if [ "$#" -gt 0 ]; then
	        cd "$(fasder -de 'printf %s' "$@")" || return 1
	    else
	        cd -
	    fi
	}
	`
}

func fzfAliases() string {
	return `
    jj () {
        local selection
        selection=$(fasder -Rdl "$@" | fzf -1 -0 --no-sort +m --height=10)
        if [[ -n "$selection" ]]; then
            echo "Selection: $selection"
            echo "$selection" | xargs -r fasder --add
            cd "$selection" || return 1
        else
            echo "No selection made"
            return 1
        fi
    }
    vv () {
      local selection
      # Get the selection from fasder and fzf
      selection=$(fasder -Rfl "$@" | fzf -1 -0 --no-sort +m --height=10)
      
      # Check if a selection was made
      if [[ -n "$selection" ]]; then
          # Ensure the editor is set and handle potential issues
          if [[ -z "$EDITOR" ]]; then
              echo "EDITOR environment variable is not set."
              return 1
          fi
          
          # Use xargs with -r to prevent running the editor if no selection
          echo "Selection: $selection"
          echo "$selection" | xargs -r fasder --add
          echo "$selection" | xargs -r "$EDITOR"
      else
          echo "No selection made."
          return 1
      fi
    }
  `
}
