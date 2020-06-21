package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	initFlags()

	if len(keepDirs) + len(purgeDirs) == 0 {
		fmt.Println("No path specified")
		return
	}

	deduper := NewDeduper()

	for keepI := range keepDirs {
		path := keepDirs[keepI]

		fmt.Println("Scanning keep directory", path)

		err := filepath.Walk(path, deduper.visitKeep)
		if err != nil {
			fmt.Printf("error walking the path %q: %v\n", path, err)
			return
		}
	}

	for purgeI := range purgeDirs {
		path := purgeDirs[purgeI]

		fmt.Println("Scanning purge directory", path)

		err := filepath.Walk(path, deduper.visitPurge)
		if err != nil {
			fmt.Printf("error walking the path %q: %v\n", path, err)
			return
		}
	}
}
