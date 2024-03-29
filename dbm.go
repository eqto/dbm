package dbm

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/eqto/dbm/stmt"
)

const (
	DriverMySQL     = `mysql`
	DriverSQLServer = `sqlserver`
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

type queryFunc func(string, ...interface{}) (*sql.Rows, error)
type selectFunc func(string, ...interface{}) ([]Resultset, error)
type execFunc func(string, ...interface{}) (sql.Result, error)

func Connect(driver, hostname string, port int, username, password, name string, opts ...Options) (*Connection, error) {
	drv, e := getDriver(driver)
	if e != nil {
		return nil, e
	}

	cn := &Connection{cfg: Config{
		DriverName: driver,
		Hostname:   hostname,
		Port:       port,
		Username:   username,
		Password:   password,
		Name:       name,
	}, drv: drv}

	if e := cn.Connect(opts...); e != nil {
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

func DeleteFrom(table string) *stmt.Delete {
	return stmt.Build().DeleteFrom(table)
}

func toFieldname(str string) string {
	fieldname := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	fieldname = matchAllCap.ReplaceAllString(fieldname, "${1}_${2}")
	return strings.ToLower(fieldname)
}
func createFieldMap(el reflect.Type) map[string]string {
	//[dbtag]fieldname
	fieldMap := make(map[string]string)
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)
		if field.Anonymous {
			subfieldMap := createFieldMap(field.Type)
			for key, val := range subfieldMap {
				fieldMap[key] = val
			}
		} else {
			dbTag := field.Tag.Get(`db`)
			if dbTag != `` {
				if dbTag != `-` {
					fieldMap[dbTag] = field.Name
				}
			} else {
				first := field.Name[0:1]
				if strings.ToUpper(first) == first {
					fieldMap[toFieldname(field.Name)] = field.Name
				}
			}
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

		switch kind {
		case reflect.Func:
			continue
		case reflect.Ptr:
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
			case reflect.Bool:
				val = rs.BoolNil(key)
			default:
				return errors.New(`unsupported ptr type: ` + key + `:` + kind.String())
			}
		default:
			switch kind {
			case reflect.String:
				val = rs.String(key)
			case reflect.Int:
				val = rs.Int(key)
			case reflect.Float64:
				val = rs.Float(key)
			case reflect.TypeOf(time.Time{}).Kind():
				val = rs.Time(key)
			case reflect.Bool:
				val = rs.Bool(key)
			default:
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
