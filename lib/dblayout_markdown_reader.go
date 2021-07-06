// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	ITEM_ID_UNKNOWN = 0
	ITEM_ID_LAYOUT  = 1
	ITEM_ID_SCHEMA  = 2
	ITEM_ID_TABLE   = 3
	ITEM_ID_FIELD   = 4
)

type ItemIdentifier int

// -----------------------------------------------------------------------------
// DbLayoutTextParser
//
// Each item points to the current field/table/schema/layout that we are parsing
// -----------------------------------------------------------------------------
type DbLayoutTextParser struct {
	LayoutPtr *DbLayout
	SchemaPtr *DbSchemaLayout
	TablePtr  *DbTableLayout
	FieldPtr  *DbFieldLayout

	// previous comment lines
	LastItemParsed ItemIdentifier
	Comment        []string
}

// -----------------------------------------------------------------------------
// NewDbLayoutTextParser
// -----------------------------------------------------------------------------
func NewDbLayoutTextParser(layoutPtr *DbLayout) DbLayoutTextParser {
	return DbLayoutTextParser{
		LayoutPtr:      layoutPtr,
		SchemaPtr:      nil,
		TablePtr:       nil,
		FieldPtr:       nil,
		LastItemParsed: ITEM_ID_UNKNOWN,
		Comment:        []string{},
	}
}

// -----------------------------------------------------------------------------
// NewDbLayoutTextParser
// -----------------------------------------------------------------------------
func (layoutParser *DbLayoutTextParser) ParseHeader(line string) {
	parts := strings.Split(line, " ")

	switch strings.Count(parts[0], "#") {
	case 1: // # database_name
		layoutParser.LayoutPtr.Name = parts[1]
		layoutParser.LayoutPtr.Type = strings.Trim(parts[2], "()")
		layoutParser.LastItemParsed = ITEM_ID_LAYOUT

	case 2: // ## schema_name
		newSchema := NewDbSchemaLayout(parts[1])
		layoutParser.SchemaPtr = &newSchema
		layoutParser.LayoutPtr.Schemas = append(layoutParser.LayoutPtr.Schemas, layoutParser.SchemaPtr)
		layoutParser.LastItemParsed = ITEM_ID_SCHEMA

	case 3: // ### table_name
		if layoutParser.SchemaPtr == nil {
			newSchema := NewDbSchemaLayout("")
			layoutParser.SchemaPtr = &newSchema
			layoutParser.LayoutPtr.Schemas = append(layoutParser.LayoutPtr.Schemas, layoutParser.SchemaPtr)
		}

		newTable := NewDbTableLayout(parts[1])
		layoutParser.TablePtr = &newTable
		layoutParser.SchemaPtr.Tables = append(layoutParser.SchemaPtr.Tables, layoutParser.TablePtr)
		layoutParser.LastItemParsed = ITEM_ID_TABLE

	default:
		fmt.Printf("Ignoring line. I don't know how to parse this: %s\n", line)
	}
}

// -----------------------------------------------------------------------------
// AssignCommentsToLastItem
// -----------------------------------------------------------------------------
func (layoutParser *DbLayoutTextParser) AssignCommentsToLastItem() {
	comment := strings.TrimSpace(strings.Join(layoutParser.Comment, " "))

	switch layoutParser.LastItemParsed {
	case ITEM_ID_LAYOUT:
		layoutParser.LayoutPtr.Comment = comment
	case ITEM_ID_SCHEMA:
		layoutParser.SchemaPtr.Comment = comment
	case ITEM_ID_TABLE:
		layoutParser.TablePtr.Comment = comment
	case ITEM_ID_FIELD:
		layoutParser.FieldPtr.Comment = comment
	default:
		if comment != "" {
			fmt.Println("ERROR: Don't know who to assign this comments to:", layoutParser.Comment)
		}
	}

	layoutParser.Comment = []string{}
}

// -----------------------------------------------------------------------------
// ParseField
//
//  - field_name [type / ....]
// -----------------------------------------------------------------------------
func (layoutParser *DbLayoutTextParser) ParseField(line string) {
	var name string
	var typeString string

	re := regexp.MustCompile(`\-\s+([^\s]+)\s*(\[[^\]]+\])?`)
	m := re.FindStringSubmatch(line)
	if len(m) == 2 {
		name = m[1]
	} else if len(m) == 3 {
		name = m[1]
		typeString = strings.Trim(strings.TrimSpace(m[2]), "[]")
	}

	field := NewDbFieldLayout(name)
	layoutParser.FieldPtr = &field
	layoutParser.TablePtr.Fields = append(layoutParser.TablePtr.Fields, layoutParser.FieldPtr)
	layoutParser.LastItemParsed = ITEM_ID_FIELD

	// TODO: parse type string when it becomes more complex, by splitting /
	field.Type = strings.Trim(typeString, "?")
	field.IsNullable = strings.HasSuffix(typeString, "?")
}

// -----------------------------------------------------------------------------
// ParseLine
// -----------------------------------------------------------------------------
func (layoutParser *DbLayoutTextParser) ParseLine(line string) {
	switch {
	case line == "":
		// skip

	case strings.HasPrefix(line, "#"):
		layoutParser.AssignCommentsToLastItem()
		layoutParser.ParseHeader(line)

	case strings.HasPrefix(line, "-"):
		layoutParser.AssignCommentsToLastItem()
		layoutParser.ParseField(line)

	default:
		layoutParser.Comment = append(layoutParser.Comment, line)
	}
}

// -----------------------------------------------------------------------------
// NewDbLayoutFromParsedString
//
// It will try to parse things with a simple algorithm, if the lines are really
// malformed, then it won't be able to do anything
// -----------------------------------------------------------------------------
func NewDbLayoutFromParsedString(text string) (*DbLayout, error) {
	layout := NewDbLayout("")
	layoutParser := NewDbLayoutTextParser(&layout)

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		layoutParser.ParseLine(line)
	}

	layoutParser.AssignCommentsToLastItem()
	layout.RebuildLookups()

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return &layout, nil
}

// -----------------------------------------------------------------------------
// NewDbLayoutFromParsedFile
// -----------------------------------------------------------------------------
func NewDbLayoutFromParsedFile(path string) (*DbLayout, error) {
	byteContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var contents string
	lpath := strings.ToLower(path)
	if strings.HasSuffix(lpath, ".md") || strings.HasSuffix(lpath, ".mardown") {
		// unescape backslashes
		buff := strings.Builder{}
		buff.Grow(len(byteContents))

		lastSkipped := false
		for _, c := range byteContents {
			if c == '\\' && !lastSkipped {
				lastSkipped = true
				continue
			}

			lastSkipped = false
			buff.WriteByte(c)
		}

		contents = buff.String()
	} else {
		contents = string(byteContents)
	}

	return NewDbLayoutFromParsedString(contents)
}
