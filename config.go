package main

import (
	"flag"
	"log"
	"regexp"
)

const (
	modeList = iota
	modeImport
)

type config struct {
	mode     int
	regex    *regexp.Regexp
	dir      string
	exclDirs strSet
}

func makeConfig() config {
	c := config{
		exclDirs: make(strSet),
	}

	imp := flag.String("import", "",
		"search for import matching the regular expression")
	flag.StringVar(&c.dir, "dir", "",
		"directory to search in (default - current directory)")
	flag.Var(c.exclDirs, "exclude-dir",
		"name of subdirectory to exclude from search")

	flag.Parse()

	c.setMode(*imp)
	c.validateDir()

	return c
}

func (c *config) setMode(imp string) {
	c.mode = modeList

	if imp != "" {
		r, err := regexp.Compile(imp)
		if err != nil {
			log.Fatalf("invalid regular expression for import: %s", err)
		}

		c.mode = modeImport
		c.regex = r
	}
}

func (c *config) validateDir() {
	if c.dir == "" {
		c.dir = getCurDir()
	} else {
		c.dir = expandDir(c.dir)
	}
}
