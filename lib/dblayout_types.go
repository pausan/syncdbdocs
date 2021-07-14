// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"errors"
	"sort"
)

type DbFieldLayout struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsUnique     bool
	IsNullable   bool
	Length       uint32
	Default      string
	Comment      string
}

type DbTableLayout struct {
	Name        string
	Comment     string
	Fields      []*DbFieldLayout
	FieldLookup map[string]*DbFieldLayout
}

type DbSchemaLayout struct {
	Name        string
	Comment     string
	Tables      []*DbTableLayout
	TableLookup map[string]*DbTableLayout
}

type DbLayout struct {
	Name         string
	Type         string // DbTypeXXX
	Comment      string
	Schemas      []*DbSchemaLayout
	SchemaLookup map[string]*DbSchemaLayout
}

// TODO: types + procedures + relationships

const (
	DbTypePostgres = "PostgreSQL"
	DbTypeMysql    = "MySQL"
	DbTypeMssql    = "MSSQL"
	DbTypeSqlite   = "SQLite"
)

const NoDbSchemaLayoutName = ""

// -----------------------------------------------------------------------------
// NewDbLayout
// -----------------------------------------------------------------------------
func NewDbLayout(name string) DbLayout {
	return DbLayout{
		Name:         name,
		Comment:      "",
		Schemas:      []*DbSchemaLayout{},
		SchemaLookup: make(map[string]*DbSchemaLayout),
	}
}

// -----------------------------------------------------------------------------
// NewDbSchemaLayout
// -----------------------------------------------------------------------------
func NewDbSchemaLayout(name string) DbSchemaLayout {
	return DbSchemaLayout{
		Name:        name,
		Comment:     "",
		Tables:      []*DbTableLayout{},
		TableLookup: make(map[string]*DbTableLayout),
	}
}

// -----------------------------------------------------------------------------
// NewDbTableLayout
// -----------------------------------------------------------------------------
func NewDbTableLayout(name string) DbTableLayout {
	return DbTableLayout{
		Name:        name,
		Comment:     "",
		Fields:      []*DbFieldLayout{},
		FieldLookup: make(map[string]*DbFieldLayout),
	}
}

// -----------------------------------------------------------------------------
// NewDbFieldLayout
// -----------------------------------------------------------------------------
func NewDbFieldLayout(name string) DbFieldLayout {
	return DbFieldLayout{
		Name:         name,
		Type:         "",
		IsPrimaryKey: false,
		IsUnique:     false,
		IsNullable:   false,
		Length:       0,
		Default:      "",
		Comment:      "",
	}
}

// -----------------------------------------------------------------------------
// GetOrCreateSchema
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) GetOrCreateSchema(schema string) *DbSchemaLayout {
	if dbSchemaLayout, ok := dbLayout.SchemaLookup[schema]; ok {
		return dbSchemaLayout
	}

	newDbSchemaLayout := NewDbSchemaLayout(schema)

	dbSchemaLayout := &newDbSchemaLayout
	dbLayout.Schemas = append(dbLayout.Schemas, dbSchemaLayout)
	dbLayout.SchemaLookup[schema] = dbSchemaLayout

	return dbSchemaLayout
}

// -----------------------------------------------------------------------------
// GetOrCreateTable
// -----------------------------------------------------------------------------
func (dbSchemaLayout *DbSchemaLayout) GetOrCreateTable(table string) *DbTableLayout {
	if dbTableLayout, ok := dbSchemaLayout.TableLookup[table]; ok {
		return dbTableLayout
	}

	newDbTableLayout := NewDbTableLayout(table)
	dbTableLayout := &newDbTableLayout
	dbSchemaLayout.Tables = append(dbSchemaLayout.Tables, dbTableLayout)
	dbSchemaLayout.TableLookup[table] = dbTableLayout
	return dbTableLayout
}

// -----------------------------------------------------------------------------
// AddField
// -----------------------------------------------------------------------------
func (dbTableLayout *DbTableLayout) AddField(field DbFieldLayout) error {
	if _, ok := dbTableLayout.FieldLookup[field.Name]; ok {
		return errors.New("Duplicate field '" + field.Name + "' on table '" + dbTableLayout.Name + "'")
	}

	dbTableLayout.Fields = append(dbTableLayout.Fields, &field)
	dbTableLayout.FieldLookup[field.Name] = &field

	return nil
}

// -----------------------------------------------------------------------------
// AddField
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) AddField(schema string, table string, field DbFieldLayout) error {
	dbSchemaLayout := dbLayout.GetOrCreateSchema(schema)
	dbTableLayout := dbSchemaLayout.GetOrCreateTable(table)
	return dbTableLayout.AddField(field)
}

// -----------------------------------------------------------------------------
// GetField
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) GetField(schema string, table string, field string) *DbFieldLayout {
	dbSchemaLayout := dbLayout.GetOrCreateSchema(schema)
	dbTableLayout := dbSchemaLayout.GetOrCreateTable(table)

	if dbTableLayout, ok := dbTableLayout.FieldLookup[field]; ok {
		return dbTableLayout
	}

	return nil
}

// -----------------------------------------------------------------------------
// GetTable
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) GetTable(schema string, table string) *DbTableLayout {
	dbSchemaLayout := dbLayout.GetOrCreateSchema(schema)
	return dbSchemaLayout.GetOrCreateTable(table)
}

// -----------------------------------------------------------------------------
// MergeFrom
//
// Merges schemas, tables and fields that exist on provided layout by preserving
// the order from the current layout.
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) MergeFrom(otherLayout *DbLayout) {
	mergedSchemas := make([]*DbSchemaLayout, 0, len(otherLayout.Schemas))

	// insert in order
	for _, schemaPtr := range dbLayout.Schemas {
		if otherSchemaPtr, ok := otherLayout.SchemaLookup[schemaPtr.Name]; ok {
			schemaPtr.MergeFrom(otherSchemaPtr, false)
			mergedSchemas = append(mergedSchemas, schemaPtr)
		}
	}

	// insert the ones only in other
	for _, otherSchemaPtr := range otherLayout.Schemas {
		if _, ok := dbLayout.SchemaLookup[otherSchemaPtr.Name]; !ok {
			mergedSchemas = append(mergedSchemas, otherSchemaPtr)
		}
	}

	dbLayout.Name = otherLayout.Name
	dbLayout.Type = otherLayout.Type
	dbLayout.Comment = otherLayout.Comment
	dbLayout.Schemas = mergedSchemas

	dbLayout.RebuildLookups()
}

// -----------------------------------------------------------------------------
// MergeFrom
//
// Merges tables and fields that exist on provided layout by preserving
// the order from the current layout.
// -----------------------------------------------------------------------------
func (dbSchemaLayout *DbSchemaLayout) MergeFrom(
	otherSchemaLayout *DbSchemaLayout,
	rebuildLookups bool,
) {
	mergedTables := make([]*DbTableLayout, 0, len(otherSchemaLayout.Tables))

	// insert in order
	for _, tablePtr := range dbSchemaLayout.Tables {
		if otherTableLayout, ok := otherSchemaLayout.TableLookup[tablePtr.Name]; ok {
			tablePtr.MergeFrom(otherTableLayout, rebuildLookups)
			mergedTables = append(mergedTables, tablePtr)
		}
	}

	// insert the ones only in other
	for _, otherTablePtr := range otherSchemaLayout.Tables {
		if _, ok := dbSchemaLayout.TableLookup[otherTablePtr.Name]; !ok {
			mergedTables = append(mergedTables, otherTablePtr)
		}
	}

	dbSchemaLayout.Name = otherSchemaLayout.Name
	dbSchemaLayout.Comment = otherSchemaLayout.Comment
	dbSchemaLayout.Tables = mergedTables

	if rebuildLookups {
		dbSchemaLayout.RebuildLookups()
	}
}

// -----------------------------------------------------------------------------
// MergeFrom
// -----------------------------------------------------------------------------
func (dbTableLayout *DbTableLayout) MergeFrom(
	otherTableLayout *DbTableLayout,
	rebuildLookups bool,
) {
	mergedFields := []*DbFieldLayout{}

	// insert fields in order that exist in both sides
	for _, fieldPtr := range dbTableLayout.Fields {
		if otherFieldPtr, ok := otherTableLayout.FieldLookup[fieldPtr.Name]; ok {
			mergedFields = append(mergedFields, otherFieldPtr)

			// in case one of the sides has a comment but not the other, leave the
			// one with the comment
			if otherFieldPtr.Comment == "" {
				otherFieldPtr.Comment = fieldPtr.Comment
			}
		}
	}

	// insert fields that are only in the other
	for _, otherFieldPtr := range otherTableLayout.Fields {
		if _, ok := dbTableLayout.FieldLookup[otherFieldPtr.Name]; !ok {
			mergedFields = append(mergedFields, otherFieldPtr)
		}
	}

	// TODO: leave deleted items with a flag so when we print them out
	//       they will still appear as deleted but with ~~ in the markdown

	dbTableLayout.Name = otherTableLayout.Name
	dbTableLayout.Comment = otherTableLayout.Comment
	dbTableLayout.Fields = mergedFields

	if rebuildLookups {
		dbTableLayout.RebuildLookups()
	}
}

// -----------------------------------------------------------------------------
// RebuildLookups
//
// Rebuild internal lookup schemas, tables and fields
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) RebuildLookups() {
	dbLayout.SchemaLookup = make(map[string]*DbSchemaLayout, len(dbLayout.Schemas))
	for _, schemaPtr := range dbLayout.Schemas {
		dbLayout.SchemaLookup[schemaPtr.Name] = schemaPtr

		schemaPtr.RebuildLookups()
	}
}

// -----------------------------------------------------------------------------
// RebuildLookups
//
// Rebuild internal lookup tables and fields
// -----------------------------------------------------------------------------
func (dbSchemaLayout *DbSchemaLayout) RebuildLookups() {
	dbSchemaLayout.TableLookup = make(map[string]*DbTableLayout, len(dbSchemaLayout.Tables))
	for _, tablePtr := range dbSchemaLayout.Tables {
		dbSchemaLayout.TableLookup[tablePtr.Name] = tablePtr

		tablePtr.RebuildLookups()
	}
}

// -----------------------------------------------------------------------------
// RebuildLookups
//
// Rebuild internal field lookups
// -----------------------------------------------------------------------------
func (dbTableLayout *DbTableLayout) RebuildLookups() {
	dbTableLayout.FieldLookup = make(map[string]*DbFieldLayout, len(dbTableLayout.Fields))
	for _, fieldPtr := range dbTableLayout.Fields {
		dbTableLayout.FieldLookup[fieldPtr.Name] = fieldPtr
	}
}

type bySchemaName []*DbSchemaLayout

func (a bySchemaName) Len() int           { return len(a) }
func (a bySchemaName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bySchemaName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type byTableName []*DbTableLayout

func (a byTableName) Len() int           { return len(a) }
func (a byTableName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTableName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type byFieldName []*DbFieldLayout

func (a byFieldName) Len() int           { return len(a) }
func (a byFieldName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFieldName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// -----------------------------------------------------------------------------
// Sort
//
// Sort inner objects in lexicographical order
// -----------------------------------------------------------------------------
func (dbLayout *DbLayout) Sort() {
	sort.Sort(bySchemaName(dbLayout.Schemas))

	for _, schemaPtr := range dbLayout.Schemas {
		schemaPtr.Sort()
	}
}

// -----------------------------------------------------------------------------
// Sort
//
// Sort inner objects in lexicographical order
// -----------------------------------------------------------------------------
func (dbSchemaLayout *DbSchemaLayout) Sort() {
	sort.Sort(byTableName(dbSchemaLayout.Tables))

	for _, tablePtr := range dbSchemaLayout.Tables {
		tablePtr.Sort()
	}
}

// -----------------------------------------------------------------------------
// Sort
//
// Sort inner objects in lexicographical order
// -----------------------------------------------------------------------------
func (dbTableLayout *DbTableLayout) Sort() {
	sort.Sort(byFieldName(dbTableLayout.Fields))
}
