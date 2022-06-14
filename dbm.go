package dbm

import (
	"database/sql"
	"errors"
	"reflect"
	"time"

	"github.com/eqto/dbm/stmt"
)

type queryFunc func(string, ...interface{}) (*sql.Rows, error)
type selectFunc func(string, ...interface{}) ([]Resultset, error)
type execFunc func(string, ...interface{}) (sql.Result, error)

//Connect ...
func Connect(driver, host string, port int, username, password, name string) (*Connection, error) {
	cn, e := newConnection(driver, host, port, username, password, name)
	if e != nil {
		return nil, e
	}
	if e := cn.Connect(); e != nil {
		return nil, e
	}
	return cn, nil
}

func Select(fields string) *stmt.SelectFields {
	return stmt.Build().Select(fields)
}

func InsertInto(table, fields string) *stmt.Insert {
	return stmt.Build().InsertInto(table, fields)
}

func Update(table string) *stmt.Update {
	return stmt.Build().Update(table)
}

func createFieldMap(el reflect.Type) map[string]string {
	//[dbtag]fieldname
	fieldMap := make(map[string]string)
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)
		dbTag := field.Tag.Get(`db`)
		if dbTag != `` {
			fieldMap[dbTag] = field.Name
		}
	}
	return fieldMap
}

func assignStruct(dest interface{}, fieldMap map[string]string, rs Resultset, typeOf reflect.Type) error {
	valOf := reflect.ValueOf(dest)
	for key, val := range fieldMap {
		field, _ := typeOf.FieldByName(val)
		valField := valOf.Elem().FieldByName(val)

		var val interface{}
		kind := field.Type.Kind()

		if kind == reflect.Ptr {
			kind = field.Type.Elem().Kind()
			switch kind {
			case reflect.String:
				val = rs.StringNil(key)
			case reflect.Int:
				val = rs.IntNil(key)
			case reflect.Float64:
				val = rs.FloatNil(key)
			case reflect.TypeOf(time.Time{}).Kind():
				val = rs.TimeNil(key)
			default:
				println(`unsupported type: ` + key + `:` + kind.String())
				return errors.New(`unsupported type: ` + key + `:` + kind.String())
			}
		} else {
			switch kind {
			case reflect.String:
				val = rs.String(key)
			case reflect.Int:
				val = rs.Int(key)
			case reflect.Float64:
				val = rs.Float(key)
			case reflect.TypeOf(time.Time{}).Kind():
				val = rs.Time(key)
			default:
				println(`unsupported type: ` + key + `:` + kind.String())
				return errors.New(`unsupported type: ` + key + `:` + kind.String())
			}
		}
		val = reflect.ValueOf(val)
		valField.Set(val.(reflect.Value))
	}
	return nil
}

func execQuery(driver Driver, fn queryFunc, query string, args ...interface{}) ([]Resultset, error) {
	rows, e := fn(query, driver.SanitizeParams(args)...)
	if e != nil {
		return nil, wrapErr(driver, e)
	}
	defer rows.Close()

	cols, e := rows.Columns()
	if e != nil {
		return nil, wrapErr(driver, e)
	}
	colTypes, e := rows.ColumnTypes()
	if e != nil {
		return nil, wrapErr(driver, e)
	}

	var results []Resultset

	for rows.Next() {
		contents, e := driver.BuildContents(colTypes)
		if e != nil {
			return nil, wrapErr(driver, e)
		}
		if e = rows.Scan(contents...); e != nil {
			rows.Close()
			return nil, wrapErr(driver, e)
		}

		rs := Resultset{}
		for key, val := range cols {
			rs[val] = contents[key]
		}
		results = append(results, rs)
	}
	rows.Close()

	return results, nil
}

func execQueryStruct(driver Driver, fn selectFunc, dest interface{}, query string, args ...interface{}) error {
	typeOf := reflect.TypeOf(dest)

	if typeOf.Kind() != reflect.Ptr {
		return errors.New(`dest is not a pointer`)
	}
	typeOf = typeOf.Elem()
	if typeOf.Kind() != reflect.Slice {
		return errors.New(`dest is not a slice`)
	}

	rs, e := fn(query, driver.SanitizeParams(args)...)
	if e != nil {
		return e
	}

	if len(rs) == 0 {
		return errors.New(errRecordNotFound)
	}
	elType := typeOf.Elem()
	fieldMap := createFieldMap(elType)

	slice := reflect.MakeSlice(typeOf, 0, len(rs))
	for _, val := range rs {
		el := reflect.New(elType)
		assignStruct(el.Interface(), fieldMap, val, elType)
		slice = reflect.Append(slice, el.Elem())
	}
	reflect.ValueOf(dest).Elem().Set(slice)
	return nil
}

func exec(driver Driver, fn execFunc, query string, args ...interface{}) (*Result, error) {
	res, e := fn(query, driver.SanitizeParams(args)...)
	if e != nil {
		e = wrapErr(driver, e)
	}
	return &Result{result: res}, e
}
