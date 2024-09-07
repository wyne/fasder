# See: https://github.com/charmbracelet/vhs
# Run with:
# vhs util/demo.tape

Output demo.gif

Set FontSize 18
Set Width 600
Set Height 600

Hide
Type 'source util/tape-helper.sh'
Enter
Type "cd"
Enter
Type "clear"
Enter
Sleep 500ms
Show

Sleep 1s

Type "d work"
Sleep 1s
Enter
Sleep 2s

Type "j work"
Sleep 1s
Enter
Sleep 2s

Type "j"
Sleep 1s
Enter
Sleep 2s

Type "f zsh"
Sleep 1s
Enter
Sleep 2s

