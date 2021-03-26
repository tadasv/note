package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type NoteDB struct {
	primaryIndex map[string]*Note
	reverseIndex *SetIndex
	// note -> a set of note ids
	linkIndex *SetIndex
	// notes that point to a the same note (reverse of linkIndex)
	backLinkIndex *SetIndex
}

func (db *NoteDB) Load(notebookRoot string) error {
	db.primaryIndex = make(map[string]*Note)
	db.reverseIndex = &SetIndex{
		data: map[string][]string{},
	}
	db.linkIndex = &SetIndex{
		data: map[string][]string{},
	}
	db.backLinkIndex = &SetIndex{
		data: map[string][]string{},
	}

	if err := filepath.Walk(notebookRoot, func(path string, info os.FileInfo, err error) error {
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

		noteId := path[len(notebookRoot)+1:]

		if _, ok := db.primaryIndex[noteId]; ok {
			return fmt.Errorf("duplicate key %q", noteId)
		}

		db.primaryIndex[noteId] = &Note{
			ID:   noteId,
			Text: string(data),
		}

		return nil
	}); err != nil {
		return err
	}

	for _, note := range db.primaryIndex {
		tokens := note.tokenize()
		for _, token := range tokens {
			if len(token) == 0 {
				continue
			}
			db.reverseIndex.Add(token, note.ID)
		}

		links := note.parseLinksToNotes()
		for _, l := range links {
			db.linkIndex.Add(note.ID, l)
			db.backLinkIndex.Add(l, note.ID)
		}
	}

	return nil
}

func (db *NoteDB) allNotes() []*Note {
	res := []*Note{}
	for _, note := range db.primaryIndex {
		res = append(res, note)
	}
	return res
}

func (db *NoteDB) findNotes(query string) []*Note {
	if len(query) == 0 {
		return db.allNotes()
	}

	// note id => number of matches
	matchCount := map[string]int{}

	queryTokens := tokenize(query)
	queryTokenMap := map[string]interface{}{}
	for _, t := range queryTokens {
		queryTokenMap[t] = true
	}

	for k, noteIds := range db.reverseIndex.data {
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

	res := []*Note{}

	if true {
		// precise match where all tokens were found in the doc
		for id, count := range matchCount {
			if count == len(queryTokens) {
				note := db.primaryIndex[id]
				res = append(res, note)
			}
		}
	} else {
		// non-precise match
		for i := len(sorted) - 1; i >= 0; i-- {
			p := sorted[i]
			note := db.primaryIndex[p.key]
			res = append(res, note)
		}
	}
	return res
}

func (db *NoteDB) findFloatingNotes() []*Note {
	res := []*Note{}

	for id, note := range db.primaryIndex {
		if _, ok := db.backLinkIndex.data[id]; !ok {
			res = append(res, note)
		}
	}

	return res
}

func printAllNotes(db *NoteDB) {
	for _, note := range db.primaryIndex {
		printNotePreview(db, note, "", true)
	}
}
