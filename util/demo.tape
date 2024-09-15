# See: https://github.com/charmbracelet/vhs
# Run with:
# vhs util/demo.tape

Output demo.gif

Set Shell "zsh"

Set FontSize 14
Set Width 500
Set Height 300
Set Padding 20

Require vim

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

Sleep 1s

Type "# reference a file once"
Sleep 1s
Enter
Sleep 2s

Type "ls abcdefg"
Sleep 1s
Enter
Sleep 2s

Type "# fasder will remember it"
Sleep 1s
Enter
Sleep 2s

Type "# find and open it with a substring"
Sleep 1s
Enter
Sleep 2s

Type "v abcd"
Sleep 1s
Enter
Sleep 2s

Type ":q"
Enter
Sleep 1s

Type "# it also works for directories"
Sleep 1s
Enter
Sleep 2s

Type "cd workspace"
Sleep 1s
Enter
Sleep 2s

Type "cd .."
Sleep 1s
Enter
Sleep 2s

Type "j work"
Sleep 1s
Enter
Sleep 2s

# End
Sleep 2s
