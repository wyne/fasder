package main

import (
	"strconv"

	flag "github.com/cornfeedhobo/pflag"
)

type sortMode int

const (
	sortFrecent sortMode = iota
	sortRank
	sortRecent
)

// A shared value preserves fasd's last-flag-wins behavior, including -rt/-tr.
type sortFlag struct {
	mode  *sortMode
	value sortMode
}

func (f sortFlag) Set(value string) error {
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	if enabled {
		*f.mode = f.value
	} else if *f.mode == f.value {
		*f.mode = sortFrecent
	}
	return nil
}

func (f sortFlag) String() string { return strconv.FormatBool(*f.mode == f.value) }
func (f sortFlag) Type() string   { return "bool" }

func registerSortFlags(flags *flag.FlagSet, mode *sortMode) {
	flags.VarP(sortFlag{mode, sortRank}, "rank", "r", "Sort by rank only, ignoring recency")
	flags.Lookup("rank").NoOptDefVal = "true"
	flags.VarP(sortFlag{mode, sortRecent}, "recent", "t", "Sort by recent access only, ignoring rank")
	flags.Lookup("recent").NoOptDefVal = "true"
}
