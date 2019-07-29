package main

import "flag"

type config struct {
	dir      string
	exclDirs strSet
}

func makeConfig() config {
	c := config{
		exclDirs: make(strSet),
	}

	flag.StringVar(&c.dir, "dir", "", "directory to search in (default - current directory)")
	flag.Var(c.exclDirs, "exclude-dir", "name of subdirectory to exclude from search")

	flag.Parse()

	c.validateDir()

	return c
}

func (c *config) validateDir() {
	if c.dir == "" {
		c.dir = getCurDir()
	} else {
		c.dir = expandDir(c.dir)
	}
}
