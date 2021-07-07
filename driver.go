package db

import (
	"database/sql"
)

type driver interface {
	connectionString() string
	kind() string
	buildContents(colTypes []*sql.ColumnType) ([]interface{}, error)
	insertQuery(tableName string, fields []string) string

	isDuplicate(msg string) bool
}
