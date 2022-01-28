package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"time"

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
				vals[idx] = new(float64)
			case `CHAR`, `VARCHAR`:
				vals[idx] = new(string)
			default:
				println(`Not supporting `, colType.DatabaseTypeName(), ` yet.`)
				vals[idx] = new([]byte)
			}
			if nullable {
				vals[idx] = &vals[idx]
			}
		case reflect.Struct:
			var val interface{}
			switch scanType.Name() {
			case `NullInt64`:
				val = new(int)
			case `NullFloat64`:
				val = new(float64)
			case `NullTime`:
				val = new(time.Time)
			default:
				println(`Not supporting struct `, scanType.Name(), ` yet.`)
			}
			if val != nil {
				vals[idx] = &val
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}
