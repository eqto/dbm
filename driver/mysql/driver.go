package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/eqto/dbm"
	"github.com/eqto/dbm/stmt"
	"github.com/go-sql-driver/mysql"
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
	case *stmt.Delete:
		return deleteStatement(s)
	}
	return ``
}

func (Driver) DataSourceName(cfg dbm.Config) string {
	return fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`,
		cfg.Username, cfg.Password,
		cfg.Hostname, cfg.Port,
		cfg.Name,
	)
}
func (Driver) IsDuplicate(e error) bool {
	if e, ok := e.(*mysql.MySQLError); ok {
		return e.Number == 1062
	}
	return false
}

func (Driver) BuildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
	vals := make([]interface{}, len(colTypes))
	for idx, colType := range colTypes {
		scanType := colType.ScanType()
		var val interface{}
		switch scanType.Kind() {
		case reflect.Int8:
			val = new(int8)
		case reflect.Uint8:
			val = new(uint8)
		case reflect.Int16:
			val = new(int16)
		case reflect.Uint16:
			val = new(uint16)
		case reflect.Int32:
			val = new(int32)
		case reflect.Uint32:
			val = new(uint32)
		case reflect.Int64:
			val = new(int64)
		case reflect.Uint64:
			val = new(uint64)
		case reflect.Float32:
			val = new(float32)
		case reflect.Float64:
			val = new(float64)
		case reflect.Slice:
			nullable, ok := colType.Nullable()
			if !ok {
				nullable = true
			}
			switch colType.DatabaseTypeName() {
			case `DECIMAL`:
				f := new(float64)
				if nullable {
					val = &f
				} else {
					val = f
				}
			case `CHAR`, `VARCHAR`, `TEXT`, `JSON`:
				s := new(string)
				if nullable {
					val = &s
				} else {
					val = s
				}
			default:
				println(`Not supporting `, colType.DatabaseTypeName(), ` yet.`)
				b := new([]byte)
				if nullable {
					val = &b
				} else {
					val = b
				}

			}
		case reflect.Struct:
			switch scanType.Name() {
			case `NullInt64`:
				i := new(int)
				val = &i
			case `NullFloat64`:
				f := new(float64)
				val = &f
			case `NullTime`:
				t := new(time.Time)
				val = &t
			case `NullString`:
				t := new(string)
				val = &t
			default:
				println(`Not supporting struct `, scanType.Name(), ` yet.`)
			}
		}
		if val == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		} else {
			vals[idx] = val
		}
	}
	return vals, nil
}
