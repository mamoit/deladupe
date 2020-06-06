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

var hmap map[string] []DedupFile

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if info.Mode() & os.ModeSymlink != 0 {
		//log.Printf("Not evaluating symlinks")
		return nil
	} else if info.IsDir() {
		//fmt.Printf("D: %q\n", path)	
	} else {
		f, err := os.Open(path)
		if err != nil {
			log.Print(err)
			return nil
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Print(err)
		}
		//fmt.Printf("F: %d %x %q\n", info.Size(), h.Sum(nil), path)
		hash := hex.EncodeToString(h.Sum(nil))
		hmap[hash] = append(hmap[hash], DedupFile{hash, path, nil})
	}
	return nil
}

func main() {
	hmap = make(map[string] []DedupFile)

	path := os.Args[1]
	err := filepath.Walk(path, visit)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}

	hmap1 := hmap

	hmap = make(map[string] []DedupFile)
	path = os.Args[2]
	err = filepath.Walk(path, visit)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}

	hmap2 := hmap

	for k, v := range hmap1 {
		_, ok := hmap2[k]
		if ok {
			for f := range hmap1[k] {
				fmt.Println(v[f].path, "exists")
			}
		}
	}
}
