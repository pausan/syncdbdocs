// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"context"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

// -----------------------------------------------------------------------------
// getMysqlDbLayout
// -----------------------------------------------------------------------------
func (conn *DbConnection) getMysqlDbLayout() (*DbLayout, error) {
	dbLayout := NewDbLayout(conn.dbName)
	dbLayout.Type = DbTypeMysql

	err := conn.fetchMysqlColumnInfo(&dbLayout)
	if err != nil {
		return nil, err
	}

	err = conn.fetchMysqlTableInfo(&dbLayout)
	if err != nil {
		return nil, err
	}

	return &dbLayout, nil
}

// -----------------------------------------------------------------------------
// fetchMysqlColumnInfo
// -----------------------------------------------------------------------------
func (conn *DbConnection) fetchMysqlColumnInfo(dbLayout *DbLayout) error {
	type MyColumnDef struct {
		TableName     string `db:"TABLE_NAME"`
		ColumnName    string `db:"COLUMN_NAME"`
		IsNullable    string `db:"IS_NULLABLE"` // YES | NO
		ColumnType    string `db:"COLUMN_TYPE"`
		MaxLength     uint32 `db:"MAX_LENGTH"`
		ColumnComment string `db:"COLUMN_COMMENT"`
		ColumnDefault string `db:"COLUMN_DEFAULT"` // Default value
	}

	dbFields := []MyColumnDef{}

	// schema in MYSQL refers to database, whereas we keep postgres definition
	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&dbFields,
		`SELECT TABLE_NAME,
		        COLUMN_NAME,
		        COALESCE(COLUMN_DEFAULT, '') as COLUMN_DEFAULT,
		        IS_NULLABLE,
		        COALESCE(CHARACTER_MAXIMUM_LENGTH, 0) as MAX_LENGTH,
		        COLUMN_TYPE,
		        COLUMN_COMMENT
		   FROM INFORMATION_SCHEMA.COLUMNS
  	  WHERE table_schema=?
  	`,
		conn.dbName,
	)
	if err != nil {
		return err
	}

	for _, dbField := range dbFields {
		field := NewDbFieldLayout(dbField.ColumnName)
		field.Type = dbField.ColumnType
		field.IsNullable = dbField.IsNullable == "YES"
		field.Comment = dbField.ColumnComment

		// types already have length in the type itself
		field.Length = 0

		// TODO: field.IsPrimaryKey
		// TODO: field.IsUnique
		// TODO: field.Default

		err := dbLayout.AddField(
			NoDbSchemaLayoutName,
			dbField.TableName,
			field,
		)

		if err != nil {
			log.Println("Ignoring error:", err)
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// fetchMysqlTableInfo
// -----------------------------------------------------------------------------
func (conn *DbConnection) fetchMysqlTableInfo(dbLayout *DbLayout) error {
	type MyTableDef struct {
		TableName string `db:"TABLE_NAME"`
		Comment   string `db:"TABLE_COMMENT"`
	}

	tableDefList := []MyTableDef{}

	// schema in MYSQL refers to database, whereas we keep postgres definition
	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&tableDefList,
		`SELECT TABLE_NAME,
		        TABLE_COMMENT

		   FROM INFORMATION_SCHEMA.TABLES
		  WHERE table_schema=?
  	`,
		conn.dbName,
	)
	if err != nil {
		return err
	}

	for _, tableDef := range tableDefList {
		table := dbLayout.GetTable(NoDbSchemaLayoutName, tableDef.TableName)
		if table != nil {
			table.Comment = tableDef.Comment
		}
	}

	return nil
}
