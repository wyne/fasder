# See: https://github.com/charmbracelet/vhs
# Run with:
# vhs util/demo.tape

Output demo.gif

Set Shell "zsh --no-rcs"

Set FontSize 16
Set Width 900
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

Type "touch file1 # Use files like normal to populate the database"
Sleep 1s
Enter
Sleep 2s

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

# Type "clear"
# Sleep 1s
# Enter
# Sleep 2s

# Type "# The same works with directories"
# Sleep 1s
# Enter
# Sleep 2s
# 
# Type "cd workspace # Normal actions will populate the database"
# Sleep 1s
# Enter
# Sleep 2s
# 
# Type "cd ../homework"
# Sleep 1s
# Enter
# Sleep 2s
# 
# Type ""
# Sleep 1s
# Enter
# Sleep 2s
# 
# 
# Type "f zsh"
# Sleep 1s
# Enter
# Sleep 2s
# 
# Type "d work"
# Sleep 1s
# Enter
# Sleep 2s
# 
# Type "j work"
# Sleep 1s
# Enter
# Sleep 2s
# 
# Type "j"
# Sleep 1s
# Enter
# Sleep 2s
# 
