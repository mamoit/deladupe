package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// Helper functions
func computeHash(path string) (string, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer f.Close()

	// calculate sha256
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

type DedupFile struct {
	hash string
	path string
}

type dedupContentTrack struct {
	filesBySize map[int64]*hashes
}
type hashes struct {
	filesByHash map[string][]DedupFile
}

// Deduper methods
type Deduper struct {
	hmap    dedupContentTrack
	minSize uint64
}

func NewDeduper(minSize uint64) *Deduper {
	hmap := make(map[int64]*hashes)
	return &Deduper{
		hmap:    dedupContentTrack{hmap},
		minSize: minSize,
	}
}

func (d *Deduper) trackDedupFile(size int64, hash string, file DedupFile) {
	_, ok := d.hmap.filesBySize[size]
	if !ok {
		h := &hashes{filesByHash: make(map[string][]DedupFile)}
		d.hmap.filesBySize[size] = h
	}
	defupFileList := d.hmap.filesBySize[size].filesByHash[hash]
	d.hmap.filesBySize[size].filesByHash[hash] = append(defupFileList, file)
}

func (d *Deduper) getDedupFiles(size int64, hash string) []DedupFile {
	files, ok := d.hmap.filesBySize[size].filesByHash[hash]
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
	if uint64(info.Size()) < d.minSize {
		return false
	}
	return true
}

func (d *Deduper) visitTarget(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if !d.shouldVisit(info) {
		return nil
	}

	// calculate sha256
	hash, err := computeHash(path)
	if err != nil {
		return err
	}

	dedup := DedupFile{
		hash: hash,
		path: path,
	}
	d.trackDedupFile(info.Size(), hash, dedup)

	return nil
}

func (d *Deduper) visitSource(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	if !d.shouldVisit(info) {
		return nil
	}

	// stop if there is no file in target with the same size
	size := info.Size()
	_, ok := d.hmap.filesBySize[size]
	if !ok {
		return nil
	}

	// calculate sha256
	hash, err := computeHash(path)
	if err != nil {
		return err
	}

	files := d.getDedupFiles(size, hash)
	if files != nil {
		fmt.Println(size, path)
		for file := range files {
			fmt.Println("- ", files[file].path)
		}
	}
	return nil
}

func main() {
	var minSize uint64
	if len(os.Args) >= 4 {
		minSize, _ = strconv.ParseUint(os.Args[3], 10, 64)
	}
	deduper := NewDeduper(minSize)

	path := os.Args[2]
	err := filepath.Walk(path, deduper.visitTarget)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}

	path = os.Args[1]
	err = filepath.Walk(path, deduper.visitSource)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
}
