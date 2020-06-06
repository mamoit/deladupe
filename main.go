package main

import (
	"os"
	"fmt"
	"log"
	"io"
	"strconv"
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
var min_size uint64

func visit_target(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	// Do not look into symlinks
	// Do not look into directories
	if info.Mode() & (os.ModeSymlink | os.ModeDir) != 0 {
		return nil
	}

	// check if file is too small
	if uint64(info.Size()) < min_size {
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

func visit_source(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	// Do not look into symlinks
	// do not look into directories
	if info.Mode() & (os.ModeSymlink | os.ModeDir) != 0 {
		return nil
	}

	// check if file is too small
	if uint64(info.Size()) < min_size {
		return nil
	}

	// stop if there is no file with the same size
	size := info.Size()
	_, ok := hmap[size]
	if !ok {
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

	_, ok = hmap[size][hash]
	if ok {
		fmt.Println(size, path)
		for file := range hmap[size][hash] {
			fmt.Println("- ", hmap[size][hash][file].path)
		}
	}
	return nil
}

func main() {
	hmap = make(map[int64]map[string] []DedupFile)

	min_size = 0
	if len(os.Args) >= 4 {
		min_size, _ = strconv.ParseUint(os.Args[3], 10, 64)
	}

	path := os.Args[2]
	err := filepath.Walk(path, visit_target)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}

	path = os.Args[1]
	err = filepath.Walk(path, visit_source)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
}
