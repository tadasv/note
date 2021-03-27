package main

import (
	"strings"

	"github.com/caneroj1/stemmer"
)

func tokenize(text string) []string {
	text = strings.ToUpper(text)

	bytes := []byte(text)

	for i := 0; i < len(bytes); i++ {
		c := bytes[i]
		if c < 'A' || c > 'Z' {
			bytes[i] = ' '
		}
	}

	words := strings.Split(string(bytes), " ")
	return stemmer.StemMultiple(words)
}
