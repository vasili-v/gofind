package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
)

const execPackageName = "main"

func walk(conf config) {
	if err := filepath.Walk(conf.dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			if _, ok := conf.exclDirs[info.Name()]; ok {
				return filepath.SkipDir
			}

			processDir(conf.dir, path)
		}

		return nil
	}); err != nil {
		log.Fatalf("failed while walking the path %q: %s", conf.dir, err)
	}
}

func processDir(root, path string) {
	p, err := build.ImportDir(path, 0)
	if err != nil {
		if _, ok := err.(*build.NoGoError); !ok {
			log.Print(err)
		}

		return
	}

	fmt.Printf("%s: %s", trimPath(path, root), p.Name)
	if p.Name == execPackageName {
		fmt.Printf("\n")
	} else {
		fmt.Printf(" (%s)\n", p.ImportPath)
	}
}
