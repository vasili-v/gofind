package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get current directory: %s", err)
	}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if _, ok := conf.exclDirs[info.Name()]; ok {
				return filepath.SkipDir
			}

			shortPath := strings.TrimPrefix(path, root)
			if shortPath != path {
				shortPath = filepath.Join(".", shortPath)
			}

			p, err := build.ImportDir(path, 0)
			if err != nil {
				if _, ok := err.(*build.NoGoError); !ok {
					log.Print(err)
				}

				return nil
			}

			if p.Name == "main" {
				fmt.Printf("%s: main\n", shortPath)
			} else {
				fmt.Printf("%s: %s (%s)\n", shortPath, p.Name, p.ImportPath)
			}
		}

		return nil
	}); err != nil {
		log.Fatalf("failed while walking the path %q: %s", root, err)
	}
}
