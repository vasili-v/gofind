package main

import (
	"flag"
	"log"
	"regexp"
)

type config struct {
	impRegex *regexp.Regexp
	txtRegex *regexp.Regexp
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
	txt := flag.String("text", "",
		"search for source code matching the regular expression")
	flag.BoolVar(&c.group, "group", false, "group search results")
	flag.StringVar(&c.dir, "dir", "",
		"directory to search in (default - current directory)")
	flag.Var(c.exclDirs, "exclude-dir",
		"name of subdirectory to exclude from search")

	flag.Parse()

	c.setRegex(*imp, *txt)
	c.validateDir()

	return c
}

func (c *config) setRegex(imp, txt string) {
	if imp != "" {
		r, err := regexp.Compile(imp)
		if err != nil {
			log.Fatalf("invalid regular expression for import: %s", err)
		}

		c.impRegex = r
	}

	if txt != "" {
		r, err := regexp.Compile(txt)
		if err != nil {
			log.Fatalf("invalid regular expression for source code: %s", err)
		}

		c.txtRegex = r
	}

	if c.group && c.txtRegex != nil {
		log.Println("warning: ignore group flag for search by sources")
	}
}

func (c *config) validateDir() {
	if c.dir == "" {
		c.dir = getCurDir()
	} else {
		c.dir = expandDir(c.dir)
	}
}
