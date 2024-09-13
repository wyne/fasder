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

Algorithm

- Run search twice. Once with last segment split by file extension and once without.
- (?) What if there's multiple file extensions?

Rules
The last search term must match the last segment of the path (this applies to both methods: splitting by / and splitting by / and .).
Search terms must match in the order they appear. For example, if the search is "tm con", "tm" must appear before "con" in the path segments.
The path segments do not have to be adjacent, meaning the search terms can match across non-consecutive segments, but they must follow the order of appearance.

## One word search

Match only against the last path segment.

### Example: `fasder score`

**Positive matches**

```
                         |---search-window-----|
- /Users/justin/workspace/scorepad-react-native
                          ^^^^^

                                              |-srch-win-|
- /Users/justin/workspace/scorepad-react-native/score.tsx
                                                ^^^^^
```

Negative matches:

```
                                               |-srch-w|
- /Users/justin/workspace/scorepad-react-native/android
                          XXXXX not in last segment

                                               |-srch-w|
- /Users/justin/workspace/scorepad-react-native/App.tsx
                          XXXXX not in last segment
```

## Two word search

Match against either:

- Last path segment: filename plus extension
- One path segment + last path segment

Path

```
      |--1--|--1---|----1----|--1---|-----1-----|2-|
A     /Users/justin/workspace/fasder/commandfile.go

      |--1--|--1---|----1----|--1---|------2-------|
B     /Users/justin/workspace/fasder/commandfile.go

```

Searches

- `fasder com go` A
- `fasder com g` A
- `fasder fa com` B
- `fasder work .go` B
- `fasder fa com` B
- `fasder command file go` X - the three terms force 'command' and 'file' to be in different segments

## Three word search

Match against either:

- Last path segment: filename plus extension
- One path segment + last path segment

Path

```
      |--1--|--12--|----12---|--12--|----23-----|3-|
A     /Users/justin/workspace/fasder/commandfile.go

      |--1--|--12--|----12---|---2--|------3-------|
B     /Users/justin/workspace/fasder/commandfile.go

```

Searches

- `fasder fa com go` A
- `fasder fasder commandfile go` A
- `fasder ju fa commandfile.go` B
- `fasder work .go` B
- `fasder fa com` B
- `fasder command file go` X - the three terms force 'command' and 'file' to be in different segments

# Future considerations

- If search term contains uppercase letters, apply case-sensitive matching
