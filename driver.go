package db

import (
	"database/sql"
)

type driver interface {
	connectionString() string
	kind() string
	buildContents(colTypes []*sql.ColumnType) ([]interface{}, error)
	insertQuery(tableName string, fields []string) string
	insertReturnID(tx *Tx, tableName string, fields []string, values []interface{}) (int, error)

	isDuplicate(msg string) bool
}
