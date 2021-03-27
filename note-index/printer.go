package main

import (
	"fmt"
	"io"
)

type Printer struct {
	w  io.Writer
	db *NoteDB

	indentCache map[int]string
}

func NewPrinter(w io.Writer, db *NoteDB) *Printer {
	return &Printer{
		w:  w,
		db: db,
		indentCache: map[int]string{
			0: makeIndent(0),
			1: makeIndent(1),
			2: makeIndent(2),
		},
	}
}

func (p *Printer) PrintNotesSummary(notes []*Note, maxDepth int) {
	for _, n := range notes {
		p.printNoteSummary(n, 0, maxDepth)
	}
}

func (p *Printer) printNoteSummary(n *Note, level, maxDepth int) {
	fmt.Fprintf(p.w, "%s %s %s\n", n.ID, p.getIndent(level), n.firstLine())

	if level+1 >= maxDepth {
		return
	}

	if links, ok := p.db.linkIndex.data[n.ID]; ok && len(links) > 0 {
		for _, linkID := range links {
			linkedNote, ok := p.db.primaryIndex[linkID]
			if !ok {
				continue
			}
			p.printNoteSummary(linkedNote, level+1, maxDepth)
		}
	}
}

func (p *Printer) getIndent(level int) string {
	indent, ok := p.indentCache[level]
	if !ok {
		indent = makeIndent(level)
		p.indentCache[level] = indent
	}
	return indent
}

func (p *Printer) PrintNoteDetails(n *Note) {
	fmt.Fprintf(p.w, "ID: %s\n", n.ID)
	fmt.Fprintf(p.w, "Local path: %s\n", CLI.NotebookRoot+"/"+n.ID)
	fmt.Fprintf(p.w, "First line: %s\n", n.firstLine())
	fmt.Fprintf(p.w, "Outgoing links:\n")
	if links, ok := p.db.linkIndex.data[n.ID]; ok {
		for _, l := range links {
			linkedNote := p.db.primaryIndex[l]
			p.printNoteSummary(linkedNote, 0, 1)
		}
	}

	fmt.Printf("Incoming links:\n")
	if links, ok := p.db.backLinkIndex.data[n.ID]; ok {
		for _, l := range links {
			linkedNote := p.db.primaryIndex[l]
			p.printNoteSummary(linkedNote, 0, 1)
		}
	}
}

func makeIndent(level int) string {
	indentSize := level * 2
	indent := make([]byte, indentSize)
	for i := 0; i < indentSize; i++ {
		indent[i] = ' '
	}
	return string(indent)
}
