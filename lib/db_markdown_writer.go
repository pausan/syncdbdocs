// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"fmt"
	"io"
)

// -----------------------------------------------------------------------------
// PrintMarkdown
//
// Print markdown document with all the information
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) PrintMarkdown(out io.Writer, lineLength int) {
	ww := NewWordWrap(lineLength, 0)
	wwFields := NewWordWrap(lineLength, 2)

	fmt.Fprintln(out, "# "+dbLayout.Name+" ("+dbLayout.Type+")")
	fmt.Fprintln(out)
	if len(dbLayout.Comment) > 0 {
		fmt.Fprintln(out, ww.Wrap(dbLayout.Comment))
		fmt.Fprintln(out)
	}

	// TODO: skip schema if database does not support it

	for _, schemaLayout := range dbLayout.Schemas {
		fmt.Fprintln(out, "## "+schemaLayout.Name)
		fmt.Fprintln(out)
		if len(schemaLayout.Comment) > 0 {
			fmt.Fprintln(out, ww.Wrap(schemaLayout.Comment))
			fmt.Fprintln(out)
		}

		for _, tableLayout := range schemaLayout.Tables {
			fmt.Fprintln(out, "### "+tableLayout.Name)
			fmt.Fprintln(out)
			if len(tableLayout.Comment) > 0 {
				fmt.Fprintln(out, ww.Wrap(tableLayout.Comment))
				fmt.Fprintln(out)
			}

			for _, field := range tableLayout.Fields {
				typeString := field.Type
				if field.IsNullable {
					typeString += "?"
				}

				fmt.Fprintf(out, "- %s [%s]\n", field.Name, typeString)
				if len(field.Comment) > 0 {
					fmt.Fprintln(out)
					fmt.Fprintln(out, wwFields.Wrap(field.Comment))
				}
				fmt.Fprintln(out)
			}
		}
	}
}
