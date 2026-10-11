package main

import (
	"strconv"
	"strings"

	flag "github.com/cornfeedhobo/pflag"
)

type selectionFlag uint

func (s *selectionFlag) Set(value string) error {
	n, err := strconv.ParseUint(value, 10, strconv.IntSize)
	if err == nil {
		*s = selectionFlag(n)
	}
	return err
}

func (s *selectionFlag) String() string { return strconv.FormatUint(uint64(*s), 10) }
func (s *selectionFlag) Type() string   { return "uint" }

func registerSelectionFlag(flags *flag.FlagSet) *uint {
	selection := selectionFlag(0)
	flags.Var(&selection, "select", "Select the nth match (-1, -2, ...); 0 disables selection")
	return (*uint)(&selection)
}

// Translate fasd's numeric shorthand into a regular pflag option. Consult the
// registered flags so command values (e.g. -e -2) are never treated as selectors.
func normalizeSelectionArgs(flags *flag.FlagSet, args []string) []string {
	result := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			return append(result, args[i:]...)
		}
		consumeNext := false
		if strings.HasPrefix(arg, "--") {
			name, _, hasValue := strings.Cut(arg[2:], "=")
			f := flags.Lookup(name)
			consumeNext = f != nil && f.NoOptDefVal == "" && !hasValue
		} else if strings.HasPrefix(arg, "-") {
			for j := 1; j < len(arg); j++ {
				if arg[j] >= '0' && arg[j] <= '9' {
					if decimalDigits(arg[j:]) {
						if j > 1 {
							result = append(result, arg[:j])
						}
						arg = "--select=" + arg[j:]
					}
					break
				}
				f := flags.ShorthandLookup(arg[j : j+1])
				if f == nil || (j+1 < len(arg) && arg[j+1] == '=') {
					break
				}
				if f.NoOptDefVal == "" {
					consumeNext = j == len(arg)-1
					break
				}
			}
		}
		result = append(result, arg)
		if consumeNext && i+1 < len(args) {
			i++
			result = append(result, args[i])
		}
	}
	return result
}

func decimalDigits(value string) bool {
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return value != ""
}
