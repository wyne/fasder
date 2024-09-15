# See: https://github.com/charmbracelet/vhs
# Run with:
# vhs util/overview.tape

Output overview.gif

Set Shell "zsh"

Set FontSize 16
Set Width 900
Set Height 300
Set Padding 20

Require vim

# ===== SETUP =====
Hide
Type "zsh --no-rcs"
Enter
Type 'source util/tape-helper.sh'
Enter
Type "cd"
Enter
Type "clear"
Enter
Sleep 500ms
Show
# ===== END SETUP =====

Sleep 1s

# Open ZSHRC
Type "v zsh # vi ~/.zshrc"
Sleep 1s
Enter
Sleep 2s

# Quit vim
Type ":q"
Enter
Sleep 1s

Type "ls file1"
Sleep 1s
Enter
Sleep 2s

Type "fasder -f fil # Search files added to the database with a substring"
Sleep 1s
Enter
Sleep 2s

Type "fasder -f -e vim fil # Open the best match with vim"
Sleep 1s
Enter
Sleep 2s

Type ":q"
Enter
Sleep 1s

Type "alias v='fasder -f -e vim' # Make an alias for faster use"
Sleep 1s
Enter
Sleep 2s

Type "v fil # Productivity boost"
Sleep 1s
Enter
Sleep 2s

# End
Sleep 2s

