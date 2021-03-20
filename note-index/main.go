package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	home                = os.Getenv("HOME")
	notebookRoot        = os.Getenv("NOTEBOOK_ROOT")
	defaultNotebookRoot = os.Getenv("HOME") + "/.notes"
)

func printFloatingNotes(primaryIndex *UniqueKeyValueIndex, reverseIndex *SetIndex) {
	for noteId, noteContent := range primaryIndex.data {
		if _, ok := reverseIndex.data["[["+noteId+"]]"]; !ok {
			firstLine := strings.TrimSpace(noteContent[:strings.Index(noteContent, "\n")])
			if len(firstLine) > 0 {
				fmt.Printf("%s %s\n", noteId, firstLine)
			} else {
				fmt.Printf("%s\n", noteId)
			}
		}
	}
}

func usage() {
	fmt.Printf(`usage: %s [command]

commands:

  rev
    prints reverse index. Each word points to a list of note ids
  
  float (default)
    prints a list of notes that are not linked to other notes.
    (The note id is note present in some other note).

`, os.Args[0])
	os.Exit(-1)
}

// buildPrimaryIndex creates a primary index from note files located at root.
// Each key in the index is note ID and value is the content.
func buildPrimaryIndex(root string) (*UniqueKeyValueIndex, error) {
	idx := &UniqueKeyValueIndex{
		data: map[string]string{},
	}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// skip .git. Assumes notebook is tracked in git
		if strings.Contains(path, ".git") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		noteId := path[len(root)+1:]
		idx.Add(noteId, string(data))

		return nil
	}); err != nil {
		return nil, err
	}

	return idx, nil
}

// buildReverseIndex creates a reverse index (word to a set of notes)
func buildReverseIndex(primaryIndex *UniqueKeyValueIndex) (*SetIndex, error) {
	idx := &SetIndex{
		data: map[string][]string{},
	}

	for noteId, noteContent := range primaryIndex.data {
		trimmedContent := strings.TrimSpace(strings.ReplaceAll(noteContent, "\n", " "))
		words := strings.Split(trimmedContent, " ")
		for _, word := range words {
			trimmedWord := strings.ToLower(strings.TrimSpace(word))
			if len(trimmedWord) == 0 {
				continue
			}

			idx.Add(trimmedWord, noteId)
		}
	}

	return idx, nil
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

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	primaryIndex, err := buildPrimaryIndex(root)
	if err != nil {
		panic(err)
	}

	reverseIndex, err := buildReverseIndex(primaryIndex)
	if err != nil {
		panic(err)
	}

	switch cmd {
	case "rev":
		for key, values := range reverseIndex.data {
			fmt.Printf("%s %s\n", key, strings.Join(values, " "))
		}
	case "float":
		fallthrough
	case "":
		printFloatingNotes(primaryIndex, reverseIndex)
	default:
		usage()
	}
}
