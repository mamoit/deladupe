package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Deduper struct {
	lock sync.Mutex

	filesBySize map[int64]*SameSized
}

type SameSized struct {
	lock sync.Mutex

	pending string

	filesByHash map[string][]string
}

// type hashMap struct {
// 	hashMap map[string][]DedupFile
// }

func NewDeduper() *Deduper {
	filesBySize := make(map[int64]*SameSized)

	return &Deduper{
		filesBySize: filesBySize,
	}
}

func (d *Deduper) getDedupFiles(size int64, hash string) []string {
	files, ok := d.filesBySize[size].filesByHash[hash]
	if !ok {
		return nil
	}
	return files
}

func (d *Deduper) shouldVisit(info os.FileInfo) bool {
	// Do not look into symlinks
	// Do not look into directories
	if info.Mode()&(os.ModeSymlink|os.ModeDir) != 0 {
		return false
	}

	// check if file is too small
	if uint64(info.Size()) < minSize {
		return false
	}
	return true
}

func (d *Deduper) visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if !d.shouldVisit(info) {
		return nil
	}

	size := info.Size()
	_, ok := d.filesBySize[size]

	// if there are no files with the same size
	if !ok {
		// Add path to the pending slot and carry on
		d.filesBySize[size] = &SameSized{
			pending:     path,
			filesByHash: make(map[string][]string),
		}
		return nil
	}

	// if there is a file with a pending hash, compute it
	if d.filesBySize[size].pending != "" {
		hash, err := computeHash(d.filesBySize[size].pending)
		if err != nil {
			return err
		}

		// if there is a pending hash, there must not be other
		// files with the same size yet
		d.filesBySize[size].pending = ""
		d.filesBySize[size].filesByHash[hash] = []string{hash}
	}

	// calculate sha256
	hash, err := computeHash(path)
	if err != nil {
		return err
	}

	// Check if there are files with the same size and same hash
	_, ok = d.filesBySize[size].filesByHash[hash]
	if !ok {
		// there is no such hash yet, add it and carry on
		d.filesBySize[size].filesByHash[hash] = []string{hash}
		return nil
	}

	// There's already a file with the same hash.
	// Add this new one to the list
	d.filesBySize[size].filesByHash[hash] = append(d.filesBySize[size].filesByHash[hash], path)

	// TODO bitwise comparison between both files?
	// Clashes using sha256 with the same sized file are be quite improbable though...

	// Delete the new one if it is targeted for deletion
	fmt.Println("#", size, hash)
	fmt.Println("-", path)
	if delete {
		//os.Remove(path)
	}
	// TODO Do not delete if file path is the same
	// TODO handle failed deletion (no delete permission for eg)
	// for file := range files {
	// 	fmt.Println("+", files[file].path)
	// }

	return nil
}
