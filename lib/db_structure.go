// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"errors"
)

type DbFieldLayout struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsUnique     bool
	IsNullable   bool
	Length       uint
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
	DbTypeSQLite   = "SQLite"
)

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
