package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoteLinkParsing(t *testing.T) {
	n := &Note{
		ID: "123",
		Text: `
		this is a test [[123]] another
		workd [[222]] something else [[333]]
		`,
	}

	links := n.parseLinksToNotes()
	assert.Equal(t, []string{"123", "222", "333"}, links)
}
