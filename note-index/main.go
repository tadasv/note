package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/kong"
)

var (
	editor              = os.Getenv("EDITOR")
	home                = os.Getenv("HOME")
	notebookRoot        = os.Getenv("NOTEBOOK_ROOT")
	defaultNotebookRoot = os.Getenv("HOME") + "/.notes"
)

var CLI struct {
	MaxDepth     int      `default:"1" help:"max number of link levels to print"`
	NotebookRoot string   `env:"NOTEBOOK_ROOT"`
	Query        QueryCmd `cmd:"query" aliases:"q" help:"search notebook"`
	Info         InfoCmd  `cmd:"info" aliases:"i" help:"show detailed information about a note"`
	List         ListCmd  `cmd:"list" aliases:"l" help:"list all notes in the notebook"`
	Check        CheckCmd `cmd:"check" help:"perform various integrity checks on the notebook"`
	New          NewCmd   `cmd:"new" help:"create a new note"`
}

type CheckCmd struct {
}

func (c *CheckCmd) Run() error {
	db := &NoteDB{}
	if err := db.Load(CLI.NotebookRoot); err != nil {
		return err
	}

	brokenLinks := db.findBrokenLinks()
	if len(brokenLinks.data) > 0 {
		fmt.Printf("Found broken links in these notes:\n")
		for noteId, links := range brokenLinks.data {
			fmt.Printf("%s  %s\n", noteId, strings.Join(links, ","))
		}
	}

	floatingNotes := db.findFloatingNotes()
	fmt.Printf("Found notes without any incoming links:\n")
	printer := NewPrinter(os.Stdout, db)
	printer.PrintNotesSummary(floatingNotes, CLI.MaxDepth)

	return nil
}

type NewCmd struct {
}

func (c *NewCmd) Run() error {
	id := time.Now().Format("20060102150405")
	fullPathToNoteFile, err := filepath.Abs(CLI.NotebookRoot + "/" + id)
	if err != nil {
		return err
	}

	fmt.Printf("Creating new note at %q\n", fullPathToNoteFile)
	if err := os.MkdirAll(CLI.NotebookRoot, 0755); err != nil {
		return err
	}

	cmd := exec.Command(editor, fullPathToNoteFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

type ListCmd struct {
}

func (c *ListCmd) Run() error {
	db := &NoteDB{}
	if err := db.Load(CLI.NotebookRoot); err != nil {
		return err
	}

	notes := db.allNotes()
	printer := NewPrinter(os.Stdout, db)
	printer.PrintNotesSummary(notes, CLI.MaxDepth)
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

	notes := db.findNotes(CLI.Query.Query)
	printer := NewPrinter(os.Stdout, db)
	printer.PrintNotesSummary(notes, CLI.MaxDepth)
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
	printer := NewPrinter(os.Stdout, db)
	printer.PrintNoteDetails(note)
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
