// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------
// EscapeFunc
// -----------------------------------------------------------------------------
type EscapeFunc func(string) string

// -----------------------------------------------------------------------------
// IdentityEscape
//
// It does nothing, it just returns what it gets
// -----------------------------------------------------------------------------
func IdentityEscape(text string) string {
	return text
}

// -----------------------------------------------------------------------------
// MarkdownEscape
//
// Escape markdown characters
// -----------------------------------------------------------------------------
func MarkdownEscape(text string) string {
	buff := strings.Builder{}
	buff.Grow(int(float32(len(text)) * 1.1))

	for _, r := range text {
		switch r {
		case '\\', '`', '{', '}', '[', ']', '<', '>', '(', ')', '#', '*', '+', '-', '_', '!', '|':
			buff.WriteRune('\\')
		}
		buff.WriteRune(r)
	}

	return buff.String()
}

// -----------------------------------------------------------------------------
// PrintMarkdown
//
// Print markdown document with all the information
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) PrintMarkdown(out io.Writer, lineLength int) {
	dbLayout.printTextWithEscapeFunction(out, lineLength, MarkdownEscape)
}

// -----------------------------------------------------------------------------
// PrintText
//
// Print a text document, which will be similar to markdown but without escaping
// anything at all.
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) PrintText(out io.Writer, lineLength int) {
	dbLayout.printTextWithEscapeFunction(out, lineLength, IdentityEscape)
}

// -----------------------------------------------------------------------------
// printTextWithEscapeFunction
//
// Helper function that will print markdown or text representation. The only
// difference is whether text needs to be escaped or not
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) printTextWithEscapeFunction(
	out io.Writer,
	lineLength int,
	escape EscapeFunc,
) {
	ww := NewWordWrap(lineLength, 0)
	wwFields := NewWordWrap(lineLength, 2)

	fmt.Fprintln(out, "# "+escape(dbLayout.Name)+" ("+dbLayout.Type+")")
	fmt.Fprintln(out)
	if len(dbLayout.Comment) > 0 {
		comment := escape(dbLayout.Comment)
		fmt.Fprintln(out, ww.Wrap(comment))
		fmt.Fprintln(out)
	}

	// TODO: skip schema if database does not support it

	// NOTE: names won't be escaped because the logic to escape them is
	//       triky and most of the times will render just fine. Markdown
	//       readability comes first.

	for _, schemaLayout := range dbLayout.Schemas {
		if schemaLayout.Name != NoDbSchemaLayoutName {
			fmt.Fprintln(out, "## "+escape(schemaLayout.Name))
			fmt.Fprintln(out)
			if len(schemaLayout.Comment) > 0 {
				comment := escape(schemaLayout.Comment)
				fmt.Fprintln(out, ww.Wrap(comment))
				fmt.Fprintln(out)
			}
		}

		for _, tableLayout := range schemaLayout.Tables {
			fmt.Fprintln(out, "### "+escape(tableLayout.Name))
			fmt.Fprintln(out)
			if len(tableLayout.Comment) > 0 {
				comment := escape(tableLayout.Comment)
				fmt.Fprintln(out, ww.Wrap(comment))
				fmt.Fprintln(out)
			}

			for _, field := range tableLayout.Fields {
				typeString := field.Type
				if field.Length > 0 {
					typeString += strconv.Itoa(int(field.Length))
				}

				if field.IsNullable {
					typeString += "?"
				}

				// TODO: add extra items with slashes so it is easy to parse
				fmt.Fprintf(out, "- %s [%s]\n", escape(field.Name), escape(typeString))
				if len(field.Comment) > 0 {
					fmt.Fprintln(out)
					comment := escape(field.Comment)
					fmt.Fprintln(out, wwFields.Wrap(comment))
				}
				fmt.Fprintln(out)
			}
		}
	}
}
