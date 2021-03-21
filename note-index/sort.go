package main

import (
	"sort"
)

type kv struct {
	key   string
	value int
}

type kvList []kv

func (p kvList) Len() int           { return len(p) }
func (p kvList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p kvList) Less(i, j int) bool { return p[i].value < p[j].value }

func sortIntMap(data map[string]int) kvList {
	p := make(kvList, len(data))

	i := 0
	for k, v := range data {
		p[i] = kv{key: k, value: v}
		i++
	}

	sort.Sort(p)
	return p
}
