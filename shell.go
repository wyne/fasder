package main

func AutoHook() string {
	return `
    if [ -n "${ZSH_VERSION-}" ]; then
        eval "$(fasder --init zsh-hook)"
    elif [ -n "${BASH_VERSION-}" ]; then
        eval "$(fasder --init bash-hook)"
    fi
    `
}

func BashHook() string {
	return `
    _fasder_prompt_func() {
        local _fasder_status=$? _fasder_command
        _fasder_command=$(HISTTIMEFORMAT= builtin history 1 |
            sed 's/^[[:space:]]*[0-9]*[[:space:]]*//')
        { eval "fasder --proc -- $(fasder --sanitize "$_fasder_command")"; } >> /dev/null 2>&1
        return "$_fasder_status"
    }
    _fasder_install_bash_hook() {
        local _fasder_prompt_entry
        case "$(declare -p PROMPT_COMMAND 2>/dev/null)" in
            "declare -a"*)
                for _fasder_prompt_entry in "${PROMPT_COMMAND[@]}"; do
                    [ "$_fasder_prompt_entry" = _fasder_prompt_func ] && return
                done
                PROMPT_COMMAND=(_fasder_prompt_func "${PROMPT_COMMAND[@]}")
                ;;
            *)
                case ";${PROMPT_COMMAND-};" in
                    *";_fasder_prompt_func;"*) ;;
                    *) PROMPT_COMMAND="_fasder_prompt_func${PROMPT_COMMAND:+;$PROMPT_COMMAND}" ;;
                esac
                ;;
        esac
    }
    _fasder_install_bash_hook
    unset -f _fasder_install_bash_hook
    `
}

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
            fasder --add "$selection"
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
          
          echo "Selection: $selection"
          fasder --add "$selection"
          "$EDITOR" "$selection"
      else
          echo "No selection made."
          return 1
      fi
    }
  `
}
