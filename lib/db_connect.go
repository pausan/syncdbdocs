// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"database/sql"
	"errors"
	"fmt"
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
