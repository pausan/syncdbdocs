// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

const (
	DriverPostgres = "pgx"
	DriverMysql    = "mysql"
	DriverMssql    = "mssql"
	DriverSqlite   = "sqlite"
)

type DbConnection struct {
	db               *sql.DB
	connectionString string
	driverType       string
	dbName           string
}

// -----------------------------------------------------------------------------
// NewDbConnection
// -----------------------------------------------------------------------------
func NewDbConnection() *DbConnection {
	conn := DbConnection{
		db:               nil,
		connectionString: "undefined",
		driverType:       "unknown",
		dbName:           "undefined",
	}

	return &conn
}

// -----------------------------------------------------------------------------
// GetDb
// -----------------------------------------------------------------------------
func (conn *DbConnection) GetDb() *sql.DB {
	return conn.db
}

// -----------------------------------------------------------------------------
// GetDriverType
// -----------------------------------------------------------------------------
func (conn *DbConnection) GetDriverType() string {
	return conn.driverType
}

// -----------------------------------------------------------------------------
// GetConnectionString
// -----------------------------------------------------------------------------
func (conn *DbConnection) GetConnectionString() string {
	return conn.connectionString
}

// -----------------------------------------------------------------------------
// DbConnect
//
// Helper method to connect to database.
// -----------------------------------------------------------------------------
func DbConnect(
	dbtype string,
	dbhost string,
	dbport uint,
	dbuser string,
	dbpass string,
	dbname string,
) (
	*DbConnection,
	error,
) {
	if dbtype == "" || dbtype == "auto" {
		return tryDbConnect(dbhost, dbport, dbuser, dbpass, dbname)
	}

	conn := NewDbConnection()
	conn.dbName = dbname

	switch dbtype {
	case "pg", "postgres", "pgx":
		conn.driverType = DriverPostgres
		conn.connectionString = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s", dbuser, dbpass, dbhost, dbport, dbname,
		)

	default:
		conn = nil
		return nil, errors.New("Invalid database type. Try with: pg | mysql | mssql | sqlite")
	}

	var err error
	conn.db, err = sql.Open(conn.driverType, conn.connectionString)
	if err != nil {
		return nil, err
	}

	return conn, err
}

// -----------------------------------------------------------------------------
// TryDbConnect
//
// Try to connect to database using given parameters with different drivers.
// When dbport is 0 it will just use the default port for each database
//
// -----------------------------------------------------------------------------
func tryDbConnect(
	dbhost string,
	dbport uint,
	dbuser string,
	dbpass string,
	dbname string,
) (
	*DbConnection,
	error,
) {
	isZeroPort := (dbport == 0)

	// try with postgres
	if isZeroPort {
		dbport = 5432
	}
	conn, err := DbConnect("pgx", dbhost, dbport, dbuser, dbpass, dbname)
	if err == nil {
		return conn, nil
	}

	// TODO: try with mysql | mssql | sqlite | ...

	return nil, errors.New("Could not connect to the database. Try specifying the database type.")
}

// -----------------------------------------------------------------------------
// Close
// -----------------------------------------------------------------------------
func (conn *DbConnection) Close() {
	if conn.db != nil {
		conn.db.Close()
		conn.db = nil
	}
}

// -----------------------------------------------------------------------------
// GetLayout
// -----------------------------------------------------------------------------
func (conn *DbConnection) GetLayout() (*DbLayout, error) {
	if conn.db == nil {
		return nil, errors.New("Not connected to any database")
	}

	switch conn.driverType {
	case DriverPostgres:
		return conn.getPostgresDbLayout()
	default:
		return nil, errors.New("Don't know how to read db layout for " + conn.driverType + " databases")
	}
}

// -----------------------------------------------------------------------------
// getPostgresDbLayout
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDbLayout() (*DbLayout, error) {
	dbLayout := NewDbLayout(conn.dbName)

	type PgFieldSchema struct {
		TableSchema            string
		TableName              string
		ColumnName             string
		IsNullable             string // YES | NO
		TypeName               string // varchar | timestamp | uuid | int2 | int8 | ...
		CharacterMaximumLength *uint32
	}

	pgFields := []PgFieldSchema{}

	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&pgFields,
		`SELECT table_schema, table_name, column_name, is_nullable, udt_name as type_name, character_maximum_length
       FROM information_schema.columns
      WHERE table_schema not in ('information_schema', 'pg_catalog')
		`,
	)
	if err != nil {
		return nil, err
	}

	for _, pgField := range pgFields {
		var length uint = 0
		if pgField.CharacterMaximumLength != nil {
			length = uint(*pgField.CharacterMaximumLength)
		}

		field := DbFieldLayout{
			Name:         pgField.ColumnName,
			Type:         pgField.TypeName,
			IsPrimaryKey: false, // TODO
			IsUnique:     false, // TODO
			IsNullable:   pgField.IsNullable == "YES",
			Length:       length,
			Default:      "", // TODO
			Comment:      "", // will be updated later
		}

		err := dbLayout.AddField(
			pgField.TableSchema,
			pgField.TableName,
			field,
		)

		if err != nil {
			log.Println("Ignoring error:", err)
		}
	}

	dbLayout.Type = DbTypePostgres

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
// getPostgresDbSchemaComments
// -----------------------------------------------------------------------------
func (conn *DbConnection) getPostgresDbSchemaComments(dbLayout *DbLayout) error {
	type TableComment struct {
		SchemaName string
		Comment    string
	}

	pgComments := []TableComment{}

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
