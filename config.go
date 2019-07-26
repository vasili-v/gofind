package main

import "flag"

type config struct {
	exclDirs strSet
}

var conf = config{
	exclDirs: make(strSet),
}

func init() {
	flag.Var(conf.exclDirs, "exclude-dir", "name of subdirectory to exclude from search")

	flag.Parse()
}
