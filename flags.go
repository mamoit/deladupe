package main

import (
	flag "github.com/spf13/pflag"
)

// Global variables
var delete bool
var minSize uint64

var keepDirs []string
var purgeDirs []string

// Parse command line flags
func initFlags() {
	flag.StringArrayVarP(&keepDirs, "keep", "k", nil, "Directories from where to keep all data.")
	flag.StringArrayVarP(&purgeDirs, "purge", "p", nil, "Directories from where to purge duplicates.")
	flag.Uint64VarP(&minSize, "minSize", "s", 1, "Minimum size of a file to be considered.")
	flag.BoolVarP(&delete, "delete", "d", false, "Delete files.")

	flag.Parse()
}
