package main

import (
	flag "github.com/spf13/pflag"
)

type arrayStringFlags []string

func (i *arrayStringFlags) String() string {
	return "my string representation"
}

func (i *arrayStringFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var delete bool

var keepDirs []string
var purgeDirs []string

var minSize uint64

func initFlags() {
	flag.StringArrayVar(&keepDirs, "keep", nil, "Directories from where to keep all data.")
	flag.StringArrayVar(&purgeDirs, "purge", nil, "Directories from where to purge duplicates.")
	flag.Uint64Var(&minSize, "minSize", 1, "Minimum size of a file to be considered.")
	flag.BoolVar(&delete, "delete", false, "Delete files.")

	flag.Parse()
}
