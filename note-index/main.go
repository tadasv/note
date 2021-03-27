package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

var (
	home                = os.Getenv("HOME")
	notebookRoot        = os.Getenv("NOTEBOOK_ROOT")
	defaultNotebookRoot = os.Getenv("HOME") + "/.notes"
)

var CLI struct {
	NotebookRoot string   `env:"NOTEBOOK_ROOT"`
	Float        FloatCmd `cmd:"float" aliases:"f"`
	Query        QueryCmd `cmd:"query" aliases:"q"`
	Info         InfoCmd  `cmd:"info" aliases:"i"`
	List         ListCmd  `cmd:"list" aliases:"l"`
}

type FloatCmd struct {
}

func (c *FloatCmd) Run() error {
	db := &NoteDB{}
	if err := db.Load(CLI.NotebookRoot); err != nil {
		return err
	}

	for _, note := range db.findFloatingNotes() {
		printNotePreview(db, note, "", false)
	}

	return nil
}

type ListCmd struct {
}

func (c *ListCmd) Run() error {
	db := &NoteDB{}
	if err := db.Load(CLI.NotebookRoot); err != nil {
		return err
	}

	for _, note := range db.allNotes() {
		printNotePreview(db, note, "", false)
	}
	return nil
}

type QueryCmd struct {
	Query string `arg:""`
}

func (c *QueryCmd) Run() error {
	db := &NoteDB{}
	if err := db.Load(CLI.NotebookRoot); err != nil {
		return err
	}

	for _, note := range db.findNotes(CLI.Query.Query) {
		printNotePreview(db, note, "", false)
	}
	return nil
}

type InfoCmd struct {
	ID string `arg:""`
}

func (c *InfoCmd) Run() error {
	db := &NoteDB{}
	if err := db.Load(CLI.NotebookRoot); err != nil {
		return err
	}

	note, ok := db.primaryIndex[CLI.Info.ID]
	if !ok {
		fmt.Printf("not found\n")
		return nil
	}
	printNote(db, note)

	return nil
}

func main() {
	ctx := kong.Parse(&CLI)
	if len(CLI.NotebookRoot) == 0 {
		CLI.NotebookRoot = defaultNotebookRoot
	}
	root, err := filepath.Abs(CLI.NotebookRoot)
	ctx.FatalIfErrorf(err)
	CLI.NotebookRoot = root
	ctx.FatalIfErrorf(ctx.Run())
}
