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

// noteMeta stores note metadata
type noteMeta struct {
	firstLine string
}

func pruneWord(word string) string {
	w := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(word)), "\n", "")
	/*
		for _, c := range ",.?!" {
			w = strings.ReplaceAll(w, string(c), "")
		}
	*/
	return w
}

func printReverseIndex(index map[string][]string) {
	for word, ids := range index {
		if len(ids) == 0 {
			continue
		}
		fmt.Printf("%s %s\n", word, strings.Join(ids, " "))
	}
}

func printFloatingNotes(index map[string][]string, noteIds []string, metas map[string]noteMeta) {
	for _, id := range noteIds {
		// TODO decide how I want to store references to other nodes. Right now it's [[noteid]]
		if _, ok := index["[["+id+"]]"]; !ok {
			if meta, ok := metas[id]; ok {
				fmt.Printf("%s %s\n", id, meta.firstLine)
			} else {
				fmt.Printf("%s\n", id)
			}
		}
	}
}

func buildIndex(root string) (map[string][]string, []string, map[string]noteMeta) {
	index := map[string][]string{}
	noteIds := []string{}
	noteMetas := map[string]noteMeta{}

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
		noteIds = append(noteIds, noteId)

		noteContent := string(data)

		noteMetas[noteId] = noteMeta{
			firstLine: noteContent[:strings.Index(noteContent, "\n")],
		}

		trimmedContent := strings.ReplaceAll(noteContent, "\n", " ")

		words := strings.Split(trimmedContent, " ")
		for _, word := range words {
			trimmedWord := pruneWord(word)
			if len(trimmedWord) == 0 {
				continue
			}

			list, ok := index[trimmedWord]
			if !ok {
				list = []string{}
			}

			exists := false

			for _, li := range list {
				if li == noteId {
					exists = true
					break
				}
			}

			if !exists {
				list = append(list, noteId)
				index[trimmedWord] = list
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return index, noteIds, noteMetas
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

	switch cmd {
	case "rev":
		index, _, _ := buildIndex(root)
		printReverseIndex(index)
	case "float":
		fallthrough
	case "":
		index, noteIds, metas := buildIndex(root)
		printFloatingNotes(index, noteIds, metas)
	default:
		usage()
	}
}
