package main

import (
	"log"
	"os"
	"path/filepath"
)

func walk(conf config) {
	s := newSearch(conf)

	if err := filepath.Walk(conf.dir,
		func(path string, info os.FileInfo, err error) error {
			if info != nil && info.IsDir() {
				if _, ok := conf.exclDirs[info.Name()]; ok {
					return filepath.SkipDir
				}

				s.processDir(path)
			}

			return nil
		},
	); err != nil {
		log.Fatalf("failed while walking the path %q: %s", conf.dir, err)
	}

	if s.postProcess != nil {
		s.postProcess()
	}
}
