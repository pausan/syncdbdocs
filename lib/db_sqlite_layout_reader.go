// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"context"

	"github.com/georgysavva/scany/sqlscan"
)

// -----------------------------------------------------------------------------
// getSqliteDbLayout
// -----------------------------------------------------------------------------
func (conn *DbConnection) getSqliteDbLayout() (*DbLayout, error) {
	dbLayout := NewDbLayout(conn.dbName)
	dbLayout.Type = DbTypeSqlite

	err := conn.fetchSqliteColumnInfo(&dbLayout)
	if err != nil {
		return nil, err
	}

	return &dbLayout, nil
}

// -----------------------------------------------------------------------------
// fetchSqliteColumnInfo
// -----------------------------------------------------------------------------
func (conn *DbConnection) fetchSqliteColumnInfo(dbLayout *DbLayout) error {
	// read: https://www.sqlite.org/schematab.html for more info
	tableNames := []string{}

	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&tableNames,
		`SELECT name AS table_name
		   FROM sqlite_master
		  WHERE type = 'table' and table_name != 'sqlite_sequence'
  	`,
	)
	if err != nil {
		return err
	}

	type SqliteColumnDef struct {
		Name         string `db:"name"`
		Type         string `db:"type"`
		NotNull      int    `db:"not_null"`
		DefaultValue string `db:"def_val"`
		PrimaryKey   int    `db:"pk"` // Default value
	}

	for _, tableName := range tableNames {
		columns := []SqliteColumnDef{}

		err := sqlscan.Select(
			ctx,
			conn.db,
			&columns,
			`SELECT name, type, [notnull] as not_null, COALESCE(dflt_value, '') as def_val, pk
			   FROM pragma_table_info(?)`,
			tableName,
		)
		if err != nil {
			return err
		}

		for _, col := range columns {
			field := NewDbFieldLayout(col.Name)
			field.Type = col.Type
			field.IsNullable = col.NotNull == 0
			field.IsPrimaryKey = col.PrimaryKey > 0
			field.Default = col.DefaultValue

			// types already have length in the type itself
			field.Length = 0
			// TODO: field.IsUnique

			err = dbLayout.AddField(
				NoDbSchemaLayoutName,
				tableName,
				field,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
