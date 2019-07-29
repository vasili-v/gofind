package main

import (
	"fmt"
	"go/build"
)

const execPackageName = "main"

type goPkg struct {
	path string
	pkg  *build.Package
}

func makeGoPkg(path string, pkg *build.Package) goPkg {
	return goPkg{
		path: path,
		pkg:  pkg,
	}
}

func (p goPkg) String() string {
	s := fmt.Sprintf("%s: %s", p.path, p.pkg.Name)
	if p.pkg.Name == execPackageName {
		return s
	}

	return s + fmt.Sprintf(" (%s)", p.pkg.ImportPath)
}
