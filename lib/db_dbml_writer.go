// Copyright (C) 2021 Pau Sanchez

package lib

import (
	"fmt"
	"io"
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
	fmt.Fprintln(out)
	if addNotes && len(dbLayout.Comment) > 0 {
		fmt.Fprintln(out, "Note: '"+dbmlEscape(dbLayout.Comment))
	}
}
