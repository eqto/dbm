package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"

	"github.com/eqto/dbm"
	"github.com/eqto/dbm/stmt"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbm.Register(`mysql`, &Driver{})
}

type Driver struct {
	dbm.Driver
}

func (Driver) Name() string {
	return `mysql`
}

func (Driver) SanitizeParams(values []interface{}) []interface{} {
	return values
}

func (Driver) StatementString(s interface{}) string {
	s = stmt.StatementOf(s)
	switch s := s.(type) {
	case *stmt.Select:
		return selectStatement(s)
	case *stmt.Insert:
		return insertStatement(s)
	case *stmt.Update:
		return updateStatement(s)
	}
	return ``
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
			nullable, ok := colType.Nullable()
			if !ok {
				nullable = true
			}
			switch colType.DatabaseTypeName() {
			case `DECIMAL`:
				if nullable {
					vals[idx] = new(sql.NullFloat64)
				} else {
					vals[idx] = new(float64)
				}
			case `VARCHAR`:
				if nullable {
					vals[idx] = new(sql.NullString)
				} else {
					vals[idx] = new(string)
				}
			default:
				println(`Not supporting `, colType.DatabaseTypeName(), ` yet.`)
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
			default:
				println(`Not supporting struct `, scanType.Name(), ` yet.`)
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}
