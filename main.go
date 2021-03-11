package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const defaultMode = 0755

var (
	home                = os.Getenv("HOME")
	editor              = os.Getenv("EDITOR")
	notebookRoot        = os.Getenv("NOTEBOOK_ROOT")
	defaultNotebookRoot = os.Getenv("HOME") + "/.notes"
)

func isosec(t time.Time) string {
	return t.Format("20060102150405")
}

func editFile(fullPath string) error {
	cmd := exec.Command(editor, fullPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func createNewNote(id string) error {
	root := notebookRoot
	if root == "" {
		root = defaultNotebookRoot
	}

	fullPathToNoteFile, err := filepath.Abs(root + "/" + id)
	if err != nil {
		return err
	}

	log.Printf("creating new note at %q", fullPathToNoteFile)
	if err := os.MkdirAll(root, defaultMode); err != nil {
		return err
	}

	return editFile(fullPathToNoteFile)
}

func main() {
	now := time.Now()
	noteId := isosec(now)
	if err := createNewNote(noteId); err != nil {
		panic(err)
	}
}
