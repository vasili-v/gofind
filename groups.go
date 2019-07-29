package main

type importGroups map[string][]goPkg

func makeImportGroups() importGroups {
	return importGroups{}
}

func (g importGroups) append(k string, p goPkg) []goPkg {
	if v, ok := g[k]; ok {
		return append(v, p)
	}

	return []goPkg{p}
}
