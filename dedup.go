package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Deduper struct {
	mux sync.Mutex

	filesBySize map[int64]*SameSized
}

type SameSized struct {
	mux sync.Mutex

	pending string

	filesByHash map[string][]ProcessedFile
}

type ProcessedFile struct {
	delete bool
	path string
}

// Constructor for Deduper
func NewDeduper() *Deduper {
	filesBySize := make(map[int64]*SameSized)

	return &Deduper{
		filesBySize: filesBySize,
	}
}

// Decide wether to visit a file or not
func (d *Deduper) shouldVisit(info os.FileInfo) bool {
	// Do not look into symlinks
	// Do not look into directories
	if info.Mode()&(os.ModeSymlink|os.ModeDir) != 0 {
		return false
	}

	// Do not conside if file is too small
	if uint64(info.Size()) < minSize {
		return false
	}
	return true
}

func deleteSymbol(toDelete bool) string {
	if !toDelete { return "+" }
	if delete && toDelete { return "-" }
	return "~"
}

func (d *Deduper) visitKeep(path string, info os.FileInfo, err error) error {
	return d.visit(path, info, err, false)
}

func (d *Deduper) visitPurge(path string, info os.FileInfo, err error) error {
	return d.visit(path, info, err, true)
}

func (d *Deduper) visit(path string, info os.FileInfo, err error, purge bool) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if !d.shouldVisit(info) {
		return nil
	}

	size := info.Size()

	d.mux.Lock()
	_, ok := d.filesBySize[size]

	// if there are no files with the same size
	if !ok {
		// Add path to the pending slot and carry on
		d.filesBySize[size] = &SameSized{
			pending:     path,
			filesByHash: make(map[string][]ProcessedFile),
		}
		d.mux.Unlock()
		return nil
	}
	d.mux.Unlock()

	d.filesBySize[size].mux.Lock()
	defer d.filesBySize[size].mux.Unlock()

	// if there is a file with a pending hash, compute it
	if d.filesBySize[size].pending != "" {
		hash, err := computeHash(d.filesBySize[size].pending)
		if err != nil {
			return err
		}

		// if there is a pending hash, there must not be other
		// files with the same size yet
		d.filesBySize[size].filesByHash[hash] = []ProcessedFile{{false, d.filesBySize[size].pending}}
		d.filesBySize[size].pending = ""
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
		d.filesBySize[size].filesByHash[hash] = []ProcessedFile{{false, path}}
		return nil
	}

	fmt.Println("#", size, hash)
	for otherI := range d.filesBySize[size].filesByHash[hash] {
		fmt.Println(
			deleteSymbol(d.filesBySize[size].filesByHash[hash][otherI].delete),
			d.filesBySize[size].filesByHash[hash][otherI].path)
	}

	// There's already a file with the same hash.
	// Add this new one to the list
	d.filesBySize[size].filesByHash[hash] = append(
		d.filesBySize[size].filesByHash[hash],
		ProcessedFile{purge, path})

	// Delete the new one if it is targeted for deletion
	fmt.Println(deleteSymbol(purge), path)
	if purge && delete {
		os.Remove(path)
	}

	return nil
}
