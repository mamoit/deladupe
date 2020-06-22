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

	delete = true
	minSize = 1
	keepDirs = []string{"tmp/keep1", "tmp/keep2"}
	purgeDirs = []string{"tmp/purge1", "tmp/purge2"}

	walk()

	if !exists("tmp/keep1/a/hi") {
		t.Error("Keep file deleted")
	}
	if !exists("tmp/keep1/hi") {
		t.Error("Keep file deleted")
	}
	if !exists("tmp/keep2/a/hi") {
		t.Error("Keep file deleted")
	}
	if !exists("tmp/keep2/hi") {
		t.Error("Keep file deleted")
	}
	if exists("tmp/purge1/a/hi") {
		t.Error("Duplicate purge file not deleted")
	}
	if exists("tmp/purge1/hi") {
		t.Error("Duplicate purge file not deleted")
	}
	if exists("tmp/purge2/a/hi") {
		t.Error("Duplicate purge file not deleted")
	}
	if exists("tmp/purge2/hi") {
		t.Error("Duplicate purge file not deleted")
	}

	if !exists("tmp/purge1/snowflake") {
		t.Error("Unique purge file deleted")
	}
	if exists("tmp/purge2/not-so-snowflake") {
		t.Error("Duplicate purge file deleted")
	}
}
