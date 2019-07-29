package main

import (
	"fmt"
	"go/build"
	"log"
)

type searchAction func(pkg goPkg)
type searchPostProcess func()

type search struct {
	conf        config
	strsCache   []string
	byImport    importGroups
	action      searchAction
	postProcess searchPostProcess
}

func newSearch(conf config) *search {
	s := &search{
		conf: conf,
	}

	s.action = s.printPkg

	if conf.impRegex != nil || conf.txtRegex != nil {
		if conf.group && conf.txtRegex == nil {
			s.action = s.groupByMatchedImports
			s.byImport = makeImportGroups()

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

	s.action(makeGoPkg(trimPath(path, s.conf.dir), p))
}

func (s *search) printPkg(pkg goPkg) {
	fmt.Println(pkg)
}

func (s *search) printMatchedPkgs(pkg goPkg) {
	strs := s.getStrsFormCache()
	defer s.updateStrsCache(strs)

	if s.conf.impRegex != nil {
		strs = s.matchImpRegex(pkg, strs)
	}

	if s.conf.txtRegex != nil {
		strs = s.matchTxtRegex(pkg, strs)
	}

	if len(strs) > 0 {
		fmt.Println(pkg)
		for _, s := range strs {
			fmt.Printf("\t%s\n", s)
		}
		fmt.Println()
	}
}

func (s *search) matchImpRegex(pkg goPkg, strs []string) []string {
	for _, imp := range pkg.pkg.Imports {
		if s.conf.impRegex.MatchString(imp) {
			strs = append(strs, imp)
		}
	}

	return strs
}

func (s *search) matchTxtRegex(pkg goPkg, strs []string) []string {
	for _, file := range pkg.pkg.GoFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.CFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.CXXFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.MFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.HFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.FFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.SFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.SwigFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	for _, file := range pkg.pkg.SwigCXXFiles {
		strs = appendMatchedStrings(s.conf.txtRegex, pkg.pkg.Dir, file, strs)
	}

	return strs
}

func (s *search) groupByMatchedImports(pkg goPkg) {
	for _, imp := range pkg.pkg.Imports {
		if s.conf.impRegex.MatchString(imp) {
			s.byImport[imp] = s.byImport.append(imp, pkg)
		}
	}
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
