#PS1='\[\033[35m\]\w > \[\033[0m\]'

PROMPT='%F{magenta}%~ >%f '

source ~/.zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

export ZSH=~/.oh-my-zsh
plugins=(
  zsh-completions
  zsh-syntax-highlighting
)

source $ZSH/oh-my-zsh.sh

## TODO: init fasder with an alternate db file
alias fasder=~/workspace/fasder/fasder

j() {
  if [ "$#" -gt 0 ]; then
    fasder_cd -d $1
  else
    cd -
  fi
}

alias a='fasder'
alias d='fasder -d'
alias f='fasder -f'
alias v='fasder -f -e vim'

echo -e "\033[35mThis is purple text\033[0m"

_fasder_preexec() {
  { eval "fasder --proc $(fasder --sanitize $2)"; } >> /dev/null 2>&1
}
autoload -Uz add-zsh-hook
add-zsh-hook preexec _fasder_preexec
