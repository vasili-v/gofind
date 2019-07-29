package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func appendMatchedStrings(
	r *regexp.Regexp,
	path, name string,
	strs []string,
) []string {
	fullPath := filepath.Join(path, name)
	f, err := os.Open(filepath.Clean(fullPath))
	if err != nil {
		log.Printf("warning: %s", err)
		return strs
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("warning %s: %s", fullPath, err)
		}
	}()

	s := bufio.NewScanner(f)
	line := 0
	for s.Scan() {
		line++
		text := s.Text()
		if r.MatchString(text) {
			strs = append(strs, fmt.Sprintf("%s:%d: %s", name, line, text))
		}
	}
	if err := s.Err(); err != nil {
		log.Printf("warning: %s: %s", fullPath, err)
	}

	return strs
}
