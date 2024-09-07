#PS1='\[\033[35m\]\w > \[\033[0m\]'

PROMPT='%F{magenta}%~ >%f '

source ~/.zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

export ZSH=~/.oh-my-zsh
plugins=(
  zsh-completions
  zsh-syntax-highlighting
)

source $ZSH/oh-my-zsh.sh

alias fasder=~/workspace/fasder/fasder

j() {
  if [ "$#" -gt 0 ]; then
    fasder_cd -d $1
  else
    cd -
  fi
}

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
alias a='fasder'
alias d='fasder -d'
alias f='fasder -f'

echo -e "\033[35mThis is purple text\033[0m"
