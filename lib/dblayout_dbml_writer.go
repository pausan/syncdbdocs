// Copyright (C) 2021 Pau Sanchez

package lib

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------
// dbmlEscape
//
// Escape a DBML string by escaping single quotes and normalizing internal spaces
// -----------------------------------------------------------------------------
func dbmlEscape(input string) string {
	normSpaces := strings.Join(strings.Fields(input), " ")
	escaped := strings.ReplaceAll(normSpaces, "'", "\\'")
	return "'" + escaped + "'"
}

// -----------------------------------------------------------------------------
// PrintDbml
//
// Print DBML document with table schemas
//
// More info:
//   - https://www.dbml.org/
//   - https://www.dbml.org/docs/
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) PrintDbml(out io.Writer, addNotes bool) {
	fmt.Fprintln(out, "Project "+dbLayout.Name+"{")
	fmt.Fprintln(out, "  database_type: "+dbmlEscape(dbLayout.Type))
	if addNotes && len(dbLayout.Comment) > 0 {
		fmt.Fprintln(out, "Note: '"+dbmlEscape(dbLayout.Comment))
	}
	fmt.Fprintln(out, "}")
	fmt.Fprintln(out)

	dbLayout.printDbmlTables(out, addNotes)
}

// -----------------------------------------------------------------------------
// printDbmlTables
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) printDbmlTables(out io.Writer, addNotes bool) {
	for _, schemaLayout := range dbLayout.Schemas {
		// TODO: there should be a way to manage schemas in postgres and print them
		//       in DBML, I guess... requires further investigation

		for _, tableLayout := range schemaLayout.Tables {
			fmt.Fprintln(out, "Table "+tableLayout.Name+"{")

			maxFieldNameLen := 10
			for _, field := range tableLayout.Fields {
				if maxFieldNameLen < len(field.Name) {
					maxFieldNameLen = len(field.Name)
				}
			}

			for _, field := range tableLayout.Fields {
				typeString := field.Type
				if field.Length > 0 {
					typeString += strconv.Itoa(int(field.Length))
				}
				if field.IsNullable {
					typeString += " [not null]"
				}

				//if addNotes && field.Comment != "" {
				//	strings.ReplaceAll(typeString, "]", ", note: '" + field.Comment + "'")
				//}

				fmt.Fprintf(out, "  %-*s %s\n", maxFieldNameLen, field.Name, typeString)
			}

			fmt.Fprintln(out, "}")
			fmt.Fprintln(out)
		}
	}

}
