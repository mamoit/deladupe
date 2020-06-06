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

func visit_source(path string, info os.FileInfo, err error) error {
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
		fmt.Println(path, "exists")
	}
	return nil
}

func main() {
	hmap = make(map[int64]map[string] []DedupFile)

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
