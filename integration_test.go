package main

import (
	"os"
	"testing"
)

// Write file with specified content
func writeFile(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// Check if a file exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func TestSimple(t *testing.T) {
	os.MkdirAll("tmp/keep1/a", 0750)
	writeFile("tmp/keep1/a/hi", "hi")
	writeFile("tmp/keep1/hi", "hi")

	os.MkdirAll("tmp/keep2/a", 0750)
	writeFile("tmp/keep2/a/hi", "hi")
	writeFile("tmp/keep2/hi", "hi")

	os.MkdirAll("tmp/purge1/a", 0750)
	writeFile("tmp/purge1/a/hi", "hi")
	writeFile("tmp/purge1/hi", "hi")
	writeFile("tmp/purge1/snowflake", "snowflake")
	os.MkdirAll("tmp/purge2/a", 0750)
	writeFile("tmp/purge2/a/hi", "hi")
	writeFile("tmp/purge2/hi", "hi")
	writeFile("tmp/purge2/not-so-snowflake", "snowflake")

	defer os.RemoveAll("tmp")

	delete = true
	minSize = 1
	keepDirs = []string{"tmp/keep1", "tmp/keep2"}
	purgeDirs = []string{"tmp/purge1", "tmp/purge2"}

	walk()

	if !exists("tmp/keep1/a/hi") || !exists("tmp/keep1/hi") || !exists("tmp/keep2/a/hi") || !exists("tmp/keep2/hi") {
		t.Error("Keep file deleted")
	}
	if exists("tmp/purge1/a/hi") || exists("tmp/purge1/hi") || exists("tmp/purge2/a/hi") || exists("tmp/purge2/hi") || exists("tmp/purge2/not-so-snowflake") {
		t.Error("Duplicate purge file not deleted")
	}
	if !exists("tmp/purge1/snowflake") {
		t.Error("Unique purge file deleted")
	}
}

func TestMinSize(t *testing.T) {
	os.MkdirAll("tmp/keep1", 0750)
	writeFile("tmp/keep1/zerobyte", "")

	os.MkdirAll("tmp/keep2/", 0750)
	writeFile("tmp/keep2/onebyte", "1")

	os.MkdirAll("tmp/purge1/", 0750)
	writeFile("tmp/purge1/zerobyte", "")
	os.MkdirAll("tmp/purge2/", 0750)
	writeFile("tmp/purge2/zerobyte", "")
	writeFile("tmp/purge2/onebyte", "1")

	defer os.RemoveAll("tmp")

	delete = true
	minSize = 1
	keepDirs = []string{"tmp/keep1", "tmp/keep2"}
	purgeDirs = []string{"tmp/purge1", "tmp/purge2"}

	walk()

	if !exists("tmp/keep1/zerobyte") || !exists("tmp/keep2/onebyte") {
		t.Error("Keep file deleted")
	}
	if !exists("tmp/purge1/zerobyte") || !exists("tmp/purge2/zerobyte") {
		t.Error("Below min size purge file deleted")
	}
	if exists("tmp/purge2/onebyte") {
		t.Error("Exactly min size purge file not deleted")
	}
}

func TestDeleteFlag(t *testing.T) {
	os.MkdirAll("tmp/keep", 0750)
	writeFile("tmp/keep/hi", "hi")

	os.MkdirAll("tmp/purge/", 0750)
	writeFile("tmp/purge/hi", "hi")

	defer os.RemoveAll("tmp")

	delete = false
	minSize = 1
	keepDirs = []string{"tmp/keep"}
	purgeDirs = []string{"tmp/purge"}

	walk()

	if !exists("tmp/keep/hi") {
		t.Error("Keep file deleted")
	}
	if !exists("tmp/purge/hi") {
		t.Error("File deleted with delete flag set to false")
	}
}

func TestSameSizeDifferentHash(t *testing.T) {
	os.MkdirAll("tmp/keep", 0750)
	writeFile("tmp/keep/3bytes", "hya")

	os.MkdirAll("tmp/purge/", 0750)
	writeFile("tmp/purge/3bytes", "hya")
	writeFile("tmp/purge/3bytes-but-different", "bye")

	defer os.RemoveAll("tmp")

	delete = true
	minSize = 1
	keepDirs = []string{"tmp/keep"}
	purgeDirs = []string{"tmp/purge"}

	walk()

	if !exists("tmp/keep/3bytes") {
		t.Error("Keep file deleted")
	}
	if !exists("tmp/purge/3bytes-but-different") {
		t.Error("File deleted with same size and unique content")
	}
	if exists("tmp/purge/3bytes") {
		t.Error("Duplicate purge file not deleted")
	}
}
