// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"context"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

// -----------------------------------------------------------------------------
// getPostgresDbLayout
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDbLayout() (*DbLayout, error) {
	dbLayout := NewDbLayout(conn.dbName)
	dbLayout.Type = DbTypePostgres

	type PgFieldSchema struct {
		TableSchema            string
		TableName              string
		ColumnName             string
		IsNullable             string // YES | NO
		TypeName               string // varchar | timestamp | uuid | int2 | int8 | ...
		CharacterMaximumLength uint32
	}

	pgFields := []PgFieldSchema{}

	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&pgFields,
		`SELECT table_schema,
		        table_name,
		        column_name,
		        is_nullable,
		        udt_name as type_name,
		        COALESCE(character_maximum_length, 0) as character_maximum_length
       FROM information_schema.columns
      WHERE table_schema not in ('information_schema', 'pg_catalog')
		`,
	)
	if err != nil {
		return nil, err
	}

	for _, pgField := range pgFields {
		field := NewDbFieldLayout(pgField.ColumnName)
		field.Type = pgField.TypeName
		field.IsNullable = pgField.IsNullable == "YES"
		field.Length = pgField.CharacterMaximumLength

		// TODO: field.IsPrimaryKey
		// TODO: field.IsUnique
		// TODO: field.Default

		err := dbLayout.AddField(
			pgField.TableSchema,
			pgField.TableName,
			field,
		)

		if err != nil {
			log.Println("Ignoring error:", err)
		}
	}

	// field.Comment will be updated here
	err = conn.getPostgresDbComments(&dbLayout)
	if err != nil {
		return nil, err
	}

	return &dbLayout, nil
}

// -----------------------------------------------------------------------------
// getPostgresDbComments
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDbComments(dbLayout *DbLayout) error {
	var err error

	err = conn.getPostgresDatabaseComment(dbLayout)
	if err != nil {
		return err
	}

	err = conn.getPostgresDbSchemaComments(dbLayout)
	if err != nil {
		return err
	}

	err = conn.getPostgresDbTableComments(dbLayout)
	if err != nil {
		return err
	}

	err = conn.getPostgresDbColumnComments(dbLayout)
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// getPostgresDatabaseComment
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDatabaseComment(dbLayout *DbLayout) error {
	var comments []string

	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&comments,
		`SELECT COALESCE(description, '') as comment
		   FROM pg_shdescription
 	LEFT JOIN pg_database on objoid = pg_database.oid
			WHERE datname = $1
		`,
		dbLayout.Name,
	)

	if err != nil {
		return err
	}

	if len(comments) >= 1 {
		dbLayout.Comment = comments[0]
	}

	return nil
}

// -----------------------------------------------------------------------------
// getPostgresDbSchemaComments
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDbSchemaComments(dbLayout *DbLayout) error {
	type SchemaComment struct {
		SchemaName string
		Comment    string
	}

	pgComments := []SchemaComment{}

	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&pgComments,
		`SELECT schema_name,
		        COALESCE(obj_description(schema_name::regnamespace, 'pg_namespace'), '') AS comment
       FROM information_schema.schemata
      WHERE schema_name NOT IN ('pg_catalog', 'information_schema')
        AND schema_name NOT LIKE 'pg_%';
		`,
	)

	if err != nil {
		return err
	}

	for _, comment := range pgComments {
		schema := dbLayout.GetOrCreateSchema(comment.SchemaName)
		if schema != nil {
			schema.Comment = comment.Comment
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// getPostgresDbTableComments
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDbTableComments(dbLayout *DbLayout) error {
	type TableComment struct {
		TableSchema string
		TableName   string
		Comment     string
	}

	pgComments := []TableComment{}

	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&pgComments,
		`SELECT table_schema,
		        table_name,
		        COALESCE(obj_description(CONCAT(table_schema, '.', table_name)::regclass, 'pg_class'), '') as comment
       FROM information_schema.tables
      WHERE table_type = 'BASE TABLE'
        AND table_schema NOT IN ('pg_catalog', 'information_schema')
		`,
	)

	if err != nil {
		return err
	}

	for _, comment := range pgComments {
		table := dbLayout.GetTable(comment.TableSchema, comment.TableName)
		if table != nil {
			table.Comment = comment.Comment
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// getPostgresDbColumnComments
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDbColumnComments(dbLayout *DbLayout) error {
	type ColumnComment struct {
		TableSchema string
		TableName   string
		ColumnName  string
		Comment     string
	}

	pgComments := []ColumnComment{}

	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&pgComments,
		`SELECT c.table_schema, c.table_name, c.column_name, pgd.description as comment
       FROM pg_catalog.pg_statio_all_tables as st
 INNER JOIN pg_catalog.pg_description pgd
         ON (pgd.objoid = st.relid)
 INNER JOIN information_schema.columns c
         ON ( pgd.objsubid=c.ordinal_position
         	AND c.table_schema=st.schemaname
         	AND c.table_name=st.relname
        )
		`,
	)

	if err != nil {
		return err
	}

	for _, comment := range pgComments {
		field := dbLayout.GetField(comment.TableSchema, comment.TableName, comment.ColumnName)
		if field != nil {
			field.Comment = comment.Comment
		}
	}

	return nil
}
