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

Type "# fasder provides aliases v, j, and others"
Sleep 1s
Enter
Sleep 2s

Type "# reference a file once"
Sleep 1s
Enter
Sleep 2s

Type "ls reallyReallyLongFileName.txt"
Sleep 1s
Enter
Sleep 2s

Type "# fasder will remember it"
Sleep 1s
Enter
Sleep 2s

Type "v really # => vim reallyReallyLongFileName.txt"
Sleep 1s
Enter
Sleep 2s

Type ":q"
Enter
Sleep 1s

Type "# it also works for directories"
Sleep 1s
Enter
Sleep 1s

Type "cd workspace"
Sleep 500ms
Enter
Sleep 1s

Type "cd .."
Enter
Sleep 500ms

Type "j work # => cd ~/workspace"
Sleep 1s
Enter
Sleep 2s

Type "v really # this also works from here"
Sleep 1s
Enter
Sleep 2s

Type ":q"
Sleep 1s
Enter
Sleep 2s

# End
Sleep 3s
