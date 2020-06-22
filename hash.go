package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Compute hash used for file comparison
func computeHash(path string) (string, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
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
