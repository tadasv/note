package main

import (
	"fmt"
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

/*
func formatLine(id string, firstLine string, indent string, maxChars int) string {
	line := fmt.Sprintf("%s%s %s", id, indent, firstLine)
	if len(line) >= maxChars {
		return line[:maxChars-3] + "..."
	}

	if len(line) < maxChars {
		toAdd := maxChars - len(line)
		// TODO
		padding := ""
		for i := 0; i < toAdd; i++ {
			padding += "."
		}
		return line + padding
	}

	return line
}
*/

func formatLine(id string, firstLine string, indent string, maxChars int) string {
	return fmt.Sprintf("%s%s %s", id, indent, firstLine)
}

func printNotePreview(db *NoteDB, n *Note, indent string, includeLinks bool) {
	fmt.Printf("%s\n", formatLine(n.ID, n.firstLine(), indent, 100))
	if includeLinks {
		links := n.parseLinksToNotes()
		if len(links) > 0 {
			for _, l := range links {
				linkedNote, ok := db.primaryIndex[l]
				if !ok {
					fmt.Printf("INVALID REFERENCE ERROR (%s), PARENT NOTE %s\n", l, n.ID)
				} else {
					printNotePreview(db, linkedNote, indent+"   ", false)
				}
			}
		}
	}
}

func printNote(db *NoteDB, n *Note) {
	fmt.Printf("ID: %s\n", n.ID)
	fmt.Printf("First line: %s\n", n.firstLine())
	fmt.Printf("Links:\n")
	if links, ok := db.linkIndex.data[n.ID]; ok {
		for _, l := range links {
			linkedNote := db.primaryIndex[l]
			fmt.Printf("  %s %s\n", linkedNote.ID, linkedNote.firstLine())
		}
	}

	fmt.Printf("Backlinks:\n")
	if links, ok := db.backLinkIndex.data[n.ID]; ok {
		for _, l := range links {
			linkedNote := db.primaryIndex[l]
			fmt.Printf("  %s %s\n", linkedNote.ID, linkedNote.firstLine())
		}
	}
}
