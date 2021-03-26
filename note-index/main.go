package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	home                = os.Getenv("HOME")
	notebookRoot        = os.Getenv("NOTEBOOK_ROOT")
	defaultNotebookRoot = os.Getenv("HOME") + "/.notes"
)

func usage() {
	fmt.Printf(`usage: %s [command]

commands:

  float
    prints a list of notes that are not linked to other notes.
    (The note id is not present in some other note).

  find
    performs a search on notes. Input argument is a regexp that will
	match against reverse index.

`, os.Args[0])
	os.Exit(-1)
}

func main() {
	root := notebookRoot
	if root == "" {
		root = defaultNotebookRoot
	}

	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		panic(err)
	}

	db := &NoteDB{}
	if err := db.Load(root); err != nil {
		panic(err)
	}

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "float":
		for _, note := range db.findFloatingNotes() {
			printNotePreview(db, note, "", false)
		}
	case "query":
		fallthrough
	case "q":
		if len(os.Args) < 2 {
			usage()
		}

		q := ""
		if len(os.Args) > 2 {
			q = os.Args[2]
		}

		for _, note := range db.findNotes(q) {
			printNotePreview(db, note, "", false)
		}
	case "h":
		fallthrough
	case "help":
		usage()
	default:
		if len(os.Args) != 2 {
			usage()
		}

		note, ok := db.primaryIndex[os.Args[1]]
		if !ok {
			fmt.Printf("not found\n")
			return
		}
		printNote(db, note)
	}
}
