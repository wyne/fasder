# See: https://github.com/charmbracelet/vhs
# Run with:
# vhs util/demo.tape

Output setup.gif

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

# Populate db
Type "echo 'Hello' > hello-file.txt"
Sleep 1s
Enter
Sleep 2s

Type "f hello # search recent files for hello"
Sleep 1s
Enter
Sleep 2s

Type "f -e vim hello # open the file in vim"
Sleep 1s
Enter
Sleep 2s

Type ":q"
Enter
Sleep 1s

Type "which v"
Sleep 1s
Enter
Sleep 2s

Type "v hello # Productivity boost"
Sleep 1s
Enter
Sleep 2s

# End
Sleep 2s

