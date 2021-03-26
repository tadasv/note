package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	home                = os.Getenv("HOME")
	notebookRoot        = os.Getenv("NOTEBOOK_ROOT")
	defaultNotebookRoot = os.Getenv("HOME") + "/.notes"
)

func getFirstLine(text string) string {
	firstLine := strings.TrimSpace(text[:strings.Index(text, "\n")])
	return firstLine
}

func printFloatingNotes(primaryIndex *UniqueKeyValueIndex, reverseIndex *SetIndex) {
	for noteId, noteContent := range primaryIndex.data {
		if _, ok := reverseIndex.data["[["+noteId+"]]"]; !ok {
			firstLine := getFirstLine(noteContent)
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
  
  float
    prints a list of notes that are not linked to other notes.
    (The note id is not present in some other note).

  find (default)
    performs a search on notes. Input argument is a regexp that will
	match against reverse index.

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
		words := tokenize(noteContent)
		// trimmedContent := strings.TrimSpace(strings.ReplaceAll(noteContent, "\n", " "))
		// words := strings.Split(trimmedContent, " ")
		for _, word := range words {
			if len(word) == 0 {
				continue
			}

			idx.Add(word, noteId)
		}
	}

	return idx, nil
}

func findNotes(primaryIndex *UniqueKeyValueIndex, reverseIndex *SetIndex, query string) {
	// note id => number of matches
	matchCount := map[string]int{}

	r, err := regexp.Compile(query)
	if err != nil {
		panic(err)
	}

	for k, noteIds := range reverseIndex.data {
		if r.MatchString(k) {
			for _, id := range noteIds {
				currentCount, ok := matchCount[id]
				if ok {
					matchCount[id] = currentCount + 1
				} else {
					matchCount[id] = 1
				}
			}
		}
	}

	sorted := sortIntMap(matchCount)

	for i := len(sorted) - 1; i >= 0; i-- {
		p := sorted[i]
		firstLine := getFirstLine(primaryIndex.data[p.key])
		fmt.Printf("%s %d %s\n", p.key, p.value, firstLine)
	}
}

func findNotes2(primaryIndex *UniqueKeyValueIndex, reverseIndex *SetIndex, query string) {
	// note id => number of matches
	matchCount := map[string]int{}

	queryTokens := tokenize(query)
	queryTokenMap := map[string]interface{}{}
	for _, t := range queryTokens {
		queryTokenMap[t] = true
	}

	for k, noteIds := range reverseIndex.data {
		if _, ok := queryTokenMap[k]; ok {
			for _, id := range noteIds {
				currentCount, ok := matchCount[id]
				if ok {
					matchCount[id] = currentCount + 1
				} else {
					matchCount[id] = 1
				}
			}
		}
	}

	sorted := sortIntMap(matchCount)

	if true {
		// precise match where all tokens were found in the doc
		for id, count := range matchCount {
			if count == len(queryTokens) {
				firstLine := getFirstLine(primaryIndex.data[id])
				fmt.Printf("%s %d %s\n", id, count, firstLine)
			}
		}
	} else {
		// non-precise match
		for i := len(sorted) - 1; i >= 0; i-- {
			p := sorted[i]
			firstLine := getFirstLine(primaryIndex.data[p.key])
			fmt.Printf("%s %d %s\n", p.key, p.value, firstLine)
		}
	}
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
		printFloatingNotes(primaryIndex, reverseIndex)
	case "":
		fallthrough
	case "find":
		if len(os.Args) != 3 {
			usage()
		}
		findNotes2(primaryIndex, reverseIndex, os.Args[2])
	default:
		usage()
	}
}
