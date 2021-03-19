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

func printFloatingNotes(index map[string][]string, noteIds []string) {
	for _, id := range noteIds {
		if _, ok := index[id]; !ok {
			fmt.Printf("%s\n", id)
		}
	}
}

func buildIndex(root string) (map[string][]string, []string) {
	index := map[string][]string{}
	noteIds := []string{}

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

		contents := strings.ReplaceAll(string(data), "\n", " ")

		words := strings.Split(contents, " ")
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

	return index, noteIds
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
		index, _ := buildIndex(root)
		printReverseIndex(index)
	case "float":
		fallthrough
	case "":
		index, noteIds := buildIndex(root)
		printFloatingNotes(index, noteIds)
	default:
		usage()
	}
}
