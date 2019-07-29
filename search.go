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
	impCache    []string
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
			s.impCache = make([]string, 0, 100)
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
	imports := s.getImpFormCache()
	defer s.updateImpCache(imports)

	for _, imp := range pkg.Imports {
		if s.conf.impRegex.MatchString(imp) {
			imports = append(imports, imp)
		}
	}

	if len(imports) > 0 {
		fmt.Println(sprintPkg(path, pkg, s.conf.dir))
		for _, s := range imports {
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

func (s *search) getImpFormCache() []string {
	return s.impCache
}

func (s *search) updateImpCache(imports []string) {
	if cap(imports) > cap(s.impCache) {
		s.impCache = imports[:0]
	}
}
