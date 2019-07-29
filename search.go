package main

import (
	"fmt"
	"go/build"
	"log"
)

const execPackageName = "main"

type searchAction func(path string, pkg *build.Package)
type searchPostProcess func()

type search struct {
	conf        config
	strsCache   []string
	byImport    map[string][]string
	action      searchAction
	postProcess searchPostProcess
}

func newSearch(conf config) *search {
	s := &search{
		conf: conf,
	}

	s.action = s.printPkg

	if conf.impRegex != nil {
		if conf.group {
			s.action = s.groupByMatchedImports
			s.byImport = map[string][]string{}

			s.postProcess = s.postGroupByMatchedImports
		} else {
			s.action = s.printMatchedPkgs
			s.strsCache = make([]string, 0, 100)
		}
	}

	return s
}

func (s *search) processDir(path string) {
	p, err := build.ImportDir(path, 0)
	if err != nil {
		if _, ok := err.(*build.NoGoError); !ok {
			log.Print(err)
		}

		return
	}

	s.action(path, p)
}

func (s *search) printPkg(path string, pkg *build.Package) {
	fmt.Println(sprintPkg(path, pkg, s.conf.dir))
}

func (s *search) printMatchedPkgs(path string, pkg *build.Package) {
	strs := s.getStrsFormCache()
	defer s.updateStrsCache(strs)

	for _, imp := range pkg.Imports {
		if s.conf.impRegex.MatchString(imp) {
			strs = append(strs, imp)
		}
	}

	if len(strs) > 0 {
		fmt.Println(sprintPkg(path, pkg, s.conf.dir))
		for _, s := range strs {
			fmt.Printf("\t%s\n", s)
		}
		fmt.Println()
	}
}

func (s *search) groupByMatchedImports(path string, pkg *build.Package) {
	var desc string
	for _, imp := range pkg.Imports {
		if s.conf.impRegex.MatchString(imp) {
			if desc == "" {
				desc = sprintPkg(path, pkg, s.conf.dir)
			}

			s.byImport[imp] = appendStringToMapValue(s.byImport, imp, desc)
		}
	}
}

func appendStringToMapValue(m map[string][]string, k, s string) []string {
	if v, ok := m[k]; ok {
		return append(v, s)
	}

	return []string{s}
}

func (s *search) postGroupByMatchedImports() {
	for k, v := range s.byImport {
		fmt.Printf("%s:\n", k)
		for _, s := range v {
			fmt.Printf("\t%s\n", s)
		}

		fmt.Println()
	}
}

func (s *search) getStrsFormCache() []string {
	return s.strsCache
}

func (s *search) updateStrsCache(strs []string) {
	if cap(strs) > cap(s.strsCache) {
		s.strsCache = strs[:0]
	}
}
