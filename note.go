package main

import (
	"regexp"
	"strings"
)

var (
	noteIdLinkRegexp = regexp.MustCompile(`\[\[(\d+)\]\]`)
)

type Note struct {
	ID   string
	Text string
}

func (n *Note) firstLine() string {
	return strings.TrimSpace(n.Text[:strings.Index(n.Text, "\n")])
}

func (n *Note) parseLinksToNotes() []string {
	matches := noteIdLinkRegexp.FindAllStringSubmatch(n.Text, -1)
	res := []string{}

	for _, match := range matches {
		res = append(res, match[1])
	}
	return res
}

func (n *Note) tokenize() []string {
	return tokenize(n.Text)
}
