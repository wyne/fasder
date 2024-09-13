# Matching

```
Search: file
Search: file ext
Search: project file
Search: username project file

Search: username project file ext
Search: username project file.ext

    /Users/username/.config/project/file.ext.ext2

[1]  Users username .config project file.ext.ext2

[2]  Users username .config project file ext ext2
```

========== Scenario 1 ==========
One search term

Example Search:

- score

Positive matches

- /Users/justin/workspace/scorepad-react-native
- /Users/justin/workspace/scorepad-react-native/score.tsx

Negative matches:

- /Users/justin/workspace/scorepad-react-native/android
- /Users/justin/workspace/scorepad-react-native/App.tsx

Approach

- Match only last path segment

========== Scenario 2 ==========
Two search terms

Example Search:

- score

Positive matches

- /Users/justin/workspace/scorepad-react-native

Negative matches:

- /Users/justin/workspace/scorepad-react-native/android

========== Scenario 3 ==========

# Future considerations

- If search term contains uppercase letters, apply case-sensitive matching
