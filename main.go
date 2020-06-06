package main

import (
	"os"
	"fmt"
	"log"
	"io"
	"path/filepath"
	"crypto/sha256"
	"encoding/hex"
)

type DedupDir struct {
	parent *DedupDir
	children []*DedupDir
	files []*DedupFile
}

type DedupFile struct {
	hash string
	path string
	parent *DedupDir
}

var hmap map[int64]map[string][]DedupFile

func visit_target(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	// Do not look into symlinks
	if info.Mode() & os.ModeSymlink != 0 {
		return nil
	}

	// do not look into directories
	if info.IsDir() {
		return nil
	}

	// open file
	f, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer f.Close()

	// calculate sha256
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Print(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))

	_, ok := hmap[info.Size()]
	if !ok {
		hmap[info.Size()] = make(map[string][]DedupFile)
	}
	hmap[info.Size()][hash] = append(hmap[info.Size()][hash], DedupFile{hash, path, nil})

	return nil
}

func main() {
	hmap = make(map[int64]map[string] []DedupFile)

	path := os.Args[1]
	err := filepath.Walk(path, visit_target)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}

	hmap1 := hmap

	hmap = make(map[int64]map[string] []DedupFile)
	path = os.Args[2]
	err = filepath.Walk(path, visit_target)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}

	hmap2 := hmap

	for size, hashes := range hmap1 {
		_, ok := hmap2[size]
		if !ok {
			continue
		}
		for hash, files := range hashes {
			_, ok := hmap2[size][hash]
			if !ok {
				continue
			}
			for file := range files {
				fmt.Println(files[file].path, "exists")
			}
		}
	}
}
