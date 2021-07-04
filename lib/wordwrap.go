// Copyright (C) 2021 Pau Sanchez
package lib

import "strings"

// -----------------------------------------------------------------------------
// WordWrap
// -----------------------------------------------------------------------------
type WordWrap struct {
	MaxLength int // controls max line length
	Indent    int // controls number of spaces each line starts with
}

// -----------------------------------------------------------------------------
// PrintMarkdown
// -----------------------------------------------------------------------------
func NewWordWrap(maxLength int, indent int) WordWrap {
	return WordWrap{
		MaxLength: maxLength,
		Indent:    indent,
	}
}

// -----------------------------------------------------------------------------
// Wrap
//
// Word wrap given text using the desired options when the class was initialized
// -----------------------------------------------------------------------------
func (ww *WordWrap) Wrap(input string) string {
	lines := make([]string, 0, 1+2*len(input)/int(ww.MaxLength))

	space := " "
	newline := "\n"

	line := []string{}
	lineLen := 0
	if ww.Indent > 0 {
		line = append(line, strings.Repeat(space, ww.Indent-1))
		lineLen = ww.Indent
	}

	words := strings.Fields(input)
	for _, word := range words {
		if lineLen+len(word) >= ww.MaxLength {
			lines = append(lines, strings.Join(line, space))

			// first item is always indentation, no need to recreate
			line = []string{}
			lineLen = 0
			if ww.Indent > 0 {
				line = append(line, strings.Repeat(space, ww.Indent-1))
				lineLen = ww.Indent
			}
		}

		line = append(line, word)
		lineLen += 1 + len(word)
	}

	lines = append(lines, strings.Join(line, space))
	return strings.Join(lines, newline)
}
