package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getCurDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get current directory: %s", err)
	}

	return dir
}

func expandDir(dir string) string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("can't expand search directory: %s", err)
	}

	return dir
}
func trimPath(path, root string) string {
	shortPath := strings.TrimPrefix(path, root)
	if shortPath == path {
		return path
	}

	return filepath.Join(".", shortPath)
}
