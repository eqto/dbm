/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2017-12-18 01:37:55
 */

package db

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

//Tx ...
type Tx struct {
	tx     *sql.Tx
	finish bool
}

//Commit ...
func (t *Tx) Commit() error {
	if t.finish {
		return errors.New(`unable to commit, transaction already finish`)
	}
	return t.tx.Commit()
}

//Rollback ...
func (t *Tx) Rollback() error {
	if t.finish {
		return errors.New(`unable to rollback, transaction already finish`)
	}
	return t.tx.Rollback()
}

//Select ...
func (t *Tx) Select(query string, params ...interface{}) ([]Resultset, error) {
	rows, e := t.tx.Query(query, params...)
	if e != nil {
		return nil, errors.New(`error executing query: ` + query)
	}
	cols, e := rows.Columns()
	colTypes, e := rows.ColumnTypes()

	var results []Resultset

	for rows.Next() {
		contents := buildContents(cols, colTypes)
		e := rows.Scan(contents...)
		if e != nil {
			println(e.Error())
			return nil, e
		}
		rs := Resultset{}
		for key, val := range cols {
			rs[val] = contents[key]
		}
		results = append(results, rs)
	}
	return results, nil
}

//SelectStruct ...
func (t *Tx) SelectStruct(dest interface{}, query string, params ...interface{}) error {
	rs, e := t.Select(query, params...)
	if e != nil {
		return e
	}
	typeOf := reflect.TypeOf(dest)
	if typeOf.Kind() != reflect.Ptr {
		return errors.New(`dest is not a pointer`)
	}
	typeOf = typeOf.Elem()
	if typeOf.Kind() != reflect.Slice {
		return errors.New(`dest is not a slice`)
	}

	slice := reflect.MakeSlice(typeOf, 0, len(rs))
	if len(rs) == 0 {
		return nil
	}
	elType := typeOf.Elem()
	fieldMap := t.createFieldMap(elType)

	for _, val := range rs {
		el := reflect.New(elType)
		t.assignStruct(el.Interface(), fieldMap, val, elType)
		slice = reflect.Append(slice, el.Elem())
	}
	reflect.ValueOf(dest).Elem().Set(slice)
	return nil
}

//Get ...
func (t *Tx) Get(query string, params ...interface{}) (Resultset, error) {
	rs, e := t.Select(query, params...)
	if e != nil || rs == nil {
		return nil, e
	}
	return rs[0], nil
}

//MustGet ...
func (t *Tx) MustGet(query string, params ...interface{}) Resultset {
	rs, e := t.Get(query, params...)
	if e != nil {
		panic(e)
	}
	return rs
}

func (t *Tx) createFieldMap(el reflect.Type) map[string]string {
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

func (t *Tx) assignStruct(dest interface{}, fieldMap map[string]string, rs Resultset, typeOf reflect.Type) error {
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
				val = rs.GetString(key)
			case reflect.Int:
				val = rs.GetInt(key)
			case reflect.Float64:
				val = rs.GetFloat(key)
			case reflect.TypeOf(time.Time{}).Kind():
				val = rs.GetTime(key)
			default:
				println(`unsupported type: ` + key + `:` + kind.String())
				return errors.New(`unsupported type: ` + key + `:` + kind.String())
			}
		} else {
			switch kind {
			case reflect.String:
				val = rs.GetStringD(key)
			case reflect.Int:
				val = rs.GetIntD(key)
			case reflect.Float64:
				val = rs.GetFloatD(key)
			case reflect.TypeOf(time.Time{}).Kind():
				val = rs.GetTimeD(key)
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

//GetStruct ...
func (t *Tx) GetStruct(dest interface{}, query string, params ...interface{}) error {
	typeOf := reflect.TypeOf(dest)
	if typeOf.Kind() != reflect.Ptr {
		return errors.New(`dest is not a pointer`)
	}

	rs, e := t.Get(query, params...)
	if e != nil {
		return e
	}
	if rs == nil || len(rs) == 0 {
		return errors.New(`record not found`)
	}

	typeOf = typeOf.Elem()
	return t.assignStruct(dest, t.createFieldMap(typeOf), rs, typeOf)
}

//Exec ...
func (t *Tx) Exec(query string, params ...interface{}) (*Result, error) {
	result, e := t.tx.Exec(query, params...)
	if e != nil {
		return nil, e
	}
	return &Result{result: result}, e
}

//MustExec ...
func (t *Tx) MustExec(query string, params ...interface{}) *Result {
	result, e := t.Exec(query, params...)
	if e != nil {
		panic(e)
	}
	return result
}

//Insert ...
func (t *Tx) Insert(tableName string, dataMap map[string]interface{}) (*Result, error) {
	var names []string
	var questionMarks []string
	var values []interface{}

	for name, value := range dataMap {
		names = append(names, name)
		values = append(values, value)
		questionMarks = append(questionMarks, `?`)
	}
	var buffer bytes.Buffer
	buffer.WriteString(`INSERT INTO `)
	buffer.WriteString(tableName)
	buffer.WriteString(`(` + strings.Join(names, `, `) + `)`)
	buffer.WriteString(` VALUES(` + strings.Join(questionMarks, `, `) + `)`)

	return t.Exec(buffer.String(), values...)
}

func buildContents(cols []string, colTypes []*sql.ColumnType) []interface{} {
	contents := make([]interface{}, len(cols))
	for i := 0; i < len(colTypes); i++ {
		colType := colTypes[i]
		switch scanType := colType.ScanType(); scanType {
		case reflect.TypeOf(sql.NullInt64{}):
			var val *uint64
			contents[i] = &val
		case reflect.TypeOf(sql.NullFloat64{}):
			var val *float64
			contents[i] = &val
		case reflect.TypeOf(sql.RawBytes{}):
			regex := getRegex()
			if regex.Match([]byte(colType.DatabaseTypeName())) {
				var val *string
				contents[i] = &val
			} else {
				var val []byte
				contents[i] = &val
			}
		case reflect.TypeOf(mysql.NullTime{}):
			var val *time.Time
			contents[i] = &val
		default:
			switch scanType.Kind() {
			case reflect.Uint64:
				fallthrough
			case reflect.Int64:
				fallthrough
			case reflect.Uint32:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Uint16:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Uint8:
				fallthrough
			case reflect.Int8:
				var val *uint64
				contents[i] = &val
			case reflect.Float64:
				fallthrough
			case reflect.Float32:
				var val *float64
				contents[i] = &val
			default:
				println(cols[i] + `:` + colType.ScanType().String())
				contents[i] = &contents[i]
			}
		}
	}
	return contents
}
