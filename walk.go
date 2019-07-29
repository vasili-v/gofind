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
	if err := filepath.Walk(conf.dir,
		func(path string, info os.FileInfo, err error) error {
			if info != nil && info.IsDir() {
				if _, ok := conf.exclDirs[info.Name()]; ok {
					return filepath.SkipDir
				}

				processDir(path, conf)
			}

			return nil
		},
	); err != nil {
		log.Fatalf("failed while walking the path %q: %s", conf.dir, err)
	}
}

func processDir(path string, conf config) {
	p, err := build.ImportDir(path, 0)
	if err != nil {
		if _, ok := err.(*build.NoGoError); !ok {
			log.Print(err)
		}

		return
	}

	switch conf.mode {
	case modeList:
		printPkg(path, p, conf.dir)

	case modeImport:
		imports := []string{}
		for _, s := range p.Imports {
			if conf.regex.MatchString(s) {
				imports = append(imports, s)
			}
		}

		if len(imports) > 0 {
			printPkg(path, p, conf.dir)
			for _, s := range imports {
				fmt.Printf("\t%s\n", s)
			}
			fmt.Println()
		}
	}
}

func printPkg(path string, pkg *build.Package, root string) {
	fmt.Printf("%s: %s", trimPath(path, root), pkg.Name)
	if pkg.Name == execPackageName {
		fmt.Printf("\n")
	} else {
		fmt.Printf(" (%s)\n", pkg.ImportPath)
	}
}
