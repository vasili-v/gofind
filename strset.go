package main

import (
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
