// Copyright (C) 2021 Pau Sanchez
package lib

import (
	"context"
	"errors"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

// NOTE: on the db writer, we need to use add or update depending on whether
//       the attribute exists
//        - sp_updateextendedproperty
//        - sp_addextendedproperty

// -----------------------------------------------------------------------------
// getMssqlDbLayout
// -----------------------------------------------------------------------------
func (conn *DbConnection) getMssqlDbLayout() (*DbLayout, error) {
	dbLayout := NewDbLayout(conn.dbName)
	dbLayout.Type = DbTypeMssql

	err := conn.fetchMssqlColumnInfo(&dbLayout)
	if err != nil {
		return nil, err
	}

	err = conn.fetchMssqlLayoutComments(&dbLayout)
	if err != nil {
		return nil, err
	}

	return &dbLayout, nil
}

// -----------------------------------------------------------------------------
// fetchMssqlColumnInfo
// -----------------------------------------------------------------------------
func (conn *DbConnection) fetchMssqlColumnInfo(dbLayout *DbLayout) error {
	type MyColumnDef struct {
		TableSchema   string `db:"TABLE_SCHEMA"`
		TableName     string `db:"TABLE_NAME"`
		ColumnName    string `db:"COLUMN_NAME"`
		IsNullable    string `db:"IS_NULLABLE"` // YES | NO
		ColumnType    string `db:"COLUMN_TYPE"`
		MaxLength     uint32 `db:"MAX_LENGTH"`
		ColumnDefault string `db:"COLUMN_DEFAULT"` // Default value
	}

	dbFields := []MyColumnDef{}

	// schema in MYSQL refers to database, whereas we keep postgres definition
	ctx := context.Background()
	err := sqlscan.Select(
		ctx,
		conn.db,
		&dbFields,
		`SELECT TABLE_SCHEMA,
						TABLE_NAME,
		        COLUMN_NAME,
		        COALESCE(COLUMN_DEFAULT, '') as COLUMN_DEFAULT,
		        IS_NULLABLE,
		        COALESCE(CHARACTER_MAXIMUM_LENGTH, 0) as MAX_LENGTH,
		        DATA_TYPE AS COLUMN_TYPE
		   FROM INFORMATION_SCHEMA.COLUMNS
  	`,
	)
	if err != nil {
		return err
	}

	for _, dbField := range dbFields {
		field := NewDbFieldLayout(dbField.ColumnName)
		field.Type = dbField.ColumnType
		field.IsNullable = dbField.IsNullable == "YES"

		// types already have length in the type itself
		field.Length = 0

		// TODO: field.IsPrimaryKey
		// TODO: field.IsUnique
		// TODO: field.Default

		err := dbLayout.AddField(
			dbField.TableSchema,
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
// fetchMssqlLayoutComments
// -----------------------------------------------------------------------------
func (conn *DbConnection) fetchMssqlLayoutComments(dbLayout *DbLayout) error {
	var err error
	dbLayout.Comment, err = conn.fetchMssqlComment(nil, nil, nil)
	if err != nil {
		return err
	}

	// According to stackoverflow we could use the following query to order to
	// extract all comments at once, but after some tests something is odd for
	// schema comments... so I prefer a slower, but correct approach
	//
	//   SELECT S.name as schema_name, O.name AS table_name, c.name as col_name, ep.value AS value
	//   FROM sys.extended_properties EP
	//   LEFT JOIN sys.all_objects O ON ep.major_id = O.object_id
	//   LEFT JOIN sys.schemas S on O.schema_id = S.schema_id
	//   LEFT JOIN sys.columns AS c ON ep.major_id = c.object_id AND ep.minor_id = c.column_id
	//   WHERE ep.name = 'MS_Description'

	for _, schemaLayout := range dbLayout.Schemas {
		schemaLayout.Comment, err = conn.fetchMssqlComment(&schemaLayout.Name, nil, nil)
		if err != nil {
			return err
		}

		for _, tableLayout := range schemaLayout.Tables {
			tableLayout.Comment, err = conn.fetchMssqlComment(&schemaLayout.Name, &tableLayout.Name, nil)
			if err != nil {
				return err
			}
			for _, field := range tableLayout.Fields {
				field.Comment, err = conn.fetchMssqlComment(&schemaLayout.Name, &tableLayout.Name, &field.Name)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// fetchMssqlLayoutComments
//
// Helper function to get a comment for an object:
//   - database comment with everything null
//   - schema with only schema as string
//   - table with schema + table set
//   - column with schema + table + column set
// -----------------------------------------------------------------------------
func (conn *DbConnection) fetchMssqlComment(schema *string, table *string, column *string) (string, error) {
	var result []string
	var err error

	ctx := context.Background()

	switch {
	// query database comment
	case schema == nil && table == nil && column == nil:
		err = sqlscan.Select(
			ctx,
			conn.db,
			&result,
			`SELECT COALESCE(value, '')
			   FROM::fn_listextendedproperty('MS_Description', NULL, NULL, NULL, NULL, NULL, NULL)
				`,
		)

	// query schema comment
	case schema != nil && table == nil && column == nil:
		err = sqlscan.Select(
			ctx,
			conn.db,
			&result,
			`SELECT COALESCE(value, '')
			   FROM::fn_listextendedproperty('MS_Description', 'schema', @p1, NULL, NULL, NULL, NULL)
				`,
			*schema,
		)

	// query table comment
	case schema != nil && table != nil && column == nil:
		err = sqlscan.Select(
			ctx,
			conn.db,
			&result,
			`SELECT COALESCE(value, '')
			   FROM::fn_listextendedproperty('MS_Description', 'schema', @p1, 'table', @p2, NULL, NULL)
				`,
			*schema,
			*table,
		)

	// query column comment
	case schema != nil && table != nil && column != nil:
		err = sqlscan.Select(
			ctx,
			conn.db,
			&result,
			`SELECT COALESCE(value, '')
			   FROM::fn_listextendedproperty('MS_Description', 'schema', @p1, 'table', @p2, 'column', @p3)
				`,
			*schema,
			*table,
			*column,
		)

	default:
		return "", errors.New("Invalid combination of parameters to retrieve comment")
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return "", err
}
