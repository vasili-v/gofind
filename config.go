package main

import (
	"flag"
	"strconv"
	"strings"
)

type strSet map[string]struct{}

func (ss strSet) String() string {
	q := make([]string, 0, len(ss))
	for s := range ss {
		q = append(q, strconv.Quote(s))
	}

	return strings.Join(q, ", ")
}

func (ss strSet) Set(s string) error {
	ss[s] = struct{}{}
	return nil
}

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
