package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"

	"github.com/eqto/go-db"
	"github.com/eqto/go-db/query"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db.Register(`mysql`, &Driver{})
}

type Driver struct {
	db.Driver
}

func (Driver) Query(stmt interface{}) string {
	stmt = query.StatementOf(stmt)
	switch stmt := stmt.(type) {
	case *query.SelectStmt:
		return querySelect(stmt)
	case *query.InsertStmt:
		// return queryInsert(stmt)
	case *query.UpdateStmt:
		// return queryUpdate(stmt)
	}
	return ``
}

func (Driver) Name() string {
	return `mysql`
}

func (Driver) DataSourceName(hostname string, port int, username, password, name string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`,
		username, password,
		hostname, port,
		name,
	)
}
func (Driver) IsDuplicate(msg string) bool {
	return regexp.MustCompile(`^Duplicate entry.*`).MatchString(msg)
}

func (Driver) BuildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
	vals := make([]interface{}, len(colTypes))
	for idx, colType := range colTypes {
		scanType := colType.ScanType()
		switch scanType.Kind() {
		case reflect.Int8:
			vals[idx] = new(int8)
		case reflect.Uint8:
			vals[idx] = new(uint8)
		case reflect.Int16:
			vals[idx] = new(int16)
		case reflect.Uint16:
			vals[idx] = new(uint16)
		case reflect.Int32:
			vals[idx] = new(int32)
		case reflect.Uint32:
			vals[idx] = new(uint32)
		case reflect.Int64:
			vals[idx] = new(int64)
		case reflect.Uint64:
			vals[idx] = new(uint64)
		case reflect.Float32:
			vals[idx] = new(float32)
		case reflect.Float64:
			vals[idx] = new(float64)
		case reflect.Slice:
			switch colType.DatabaseTypeName() {
			case `DECIMAL`:
				if null, ok := colType.Nullable(); null || ok {
					vals[idx] = new(sql.NullFloat64)
				} else {
					vals[idx] = new(float64)
				}
			default:
				vals[idx] = new([]byte)
			}
		case reflect.Struct:
			switch scanType.Name() {
			case `NullInt64`:
				vals[idx] = new(sql.NullInt64)
			case `NullFloat64`:
				vals[idx] = new(sql.NullFloat64)
			case `NullTime`:
				vals[idx] = new(sql.NullTime)
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}
