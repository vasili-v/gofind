package main

import (
	"flag"
	"log"
	"regexp"
)

type config struct {
	regex    *regexp.Regexp
	group    bool
	dir      string
	exclDirs strSet
}

func makeConfig() config {
	c := config{
		exclDirs: make(strSet),
	}

	imp := flag.String("import", "",
		"search for import matching the regular expression")
	flag.BoolVar(&c.group, "group", false, "group search results")
	flag.StringVar(&c.dir, "dir", "",
		"directory to search in (default - current directory)")
	flag.Var(c.exclDirs, "exclude-dir",
		"name of subdirectory to exclude from search")

	flag.Parse()

	c.setRegex(*imp)
	c.validateDir()

	return c
}

func (c *config) setRegex(imp string) {
	if imp != "" {
		r, err := regexp.Compile(imp)
		if err != nil {
			log.Fatalf("invalid regular expression for import: %s", err)
		}

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
