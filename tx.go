package db

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"
)

//Tx ...
type Tx struct {
	db     *sql.DB
	tx     *sql.Tx
	finish bool
}

//MustRecover ...
func (t *Tx) MustRecover() {
	if r := recover(); r != nil {
		t.Rollback()
		panic(r)
	}
	t.Commit()
}

//Recover ...
func (t *Tx) Recover() {
	if r := recover(); r != nil {
		t.Rollback()
	} else {
		t.Commit()
	}
}

//Commit ...
func (t *Tx) Commit() error {
	if t.tx == nil || t.finish {
		return nil
	}
	t.finish = true
	return t.tx.Commit()
}

//Rollback ...
func (t *Tx) Rollback() error {
	if t.tx == nil || t.finish {
		return nil
	}
	t.finish = true
	return t.tx.Rollback()
}

//MustInsert ...
func (t *Tx) MustInsert(tableName string, dataMap map[string]interface{}) *Result {
	if rs, e := t.Insert(tableName, dataMap); e != nil {
		panic(e)
	} else {
		return rs
	}
}

//MustSelect ...
func (t *Tx) MustSelect(query string, params ...interface{}) []Resultset {
	if rs, e := t.Select(query, params...); e != nil {
		panic(e)
	} else {
		return rs
	}
}

//Select ...
func (t *Tx) Select(query string, params ...interface{}) ([]Resultset, error) {
	var rows *sql.Rows
	var e error
	if t.tx == nil {
		rows, e = t.db.Query(query, params...)
	} else {
		rows, e = t.tx.Query(query, params...)
	}
	if e != nil {
		return nil, wrapErr(e)
	}
	defer rows.Close()

	cols, e := rows.Columns()
	if e != nil {
		return nil, wrapErr(e)
	}
	colTypes, e := rows.ColumnTypes()
	if e != nil {
		return nil, wrapErr(e)
	}

	var results []Resultset

	for rows.Next() {
		contents := buildContents(cols, colTypes)
		e := rows.Scan(contents...)
		if e != nil {
			println(e.Error())
			rows.Close()
			return nil, wrapErr(e)
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

//SelectStruct ...
func (t *Tx) SelectStruct(dest interface{}, query string, params ...interface{}) error {
	typeOf := reflect.TypeOf(dest)

	if typeOf.Kind() != reflect.Ptr {
		return errors.New(`dest is not a pointer`)
	}
	typeOf = typeOf.Elem()
	if typeOf.Kind() != reflect.Slice {
		return errors.New(`dest is not a slice`)
	}

	rs, e := t.Select(query, params...)
	if e != nil {
		return wrapErr(e)
	}

	if len(rs) == 0 {
		return nil
	}
	elType := typeOf.Elem()
	fieldMap := t.createFieldMap(elType)

	slice := reflect.MakeSlice(typeOf, 0, len(rs))
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
		return nil, wrapErr(e)
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
		return wrapMsgErr(`record not found`)
	}

	typeOf = typeOf.Elem()
	return t.assignStruct(dest, t.createFieldMap(typeOf), rs, typeOf)
}

//Exec ...
func (t *Tx) Exec(query string, params ...interface{}) (*Result, error) {
	var res sql.Result
	var e error
	if t.tx == nil {
		res, e = t.db.Exec(query, params...)
	} else {
		res, e = t.tx.Exec(query, params...)
	}
	if e != nil {
		return nil, wrapErr(e)
	}

	return &Result{result: res}, wrapErr(e)
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
		switch colType.ScanType().Kind() {
		case reflect.String:
			var val *string
			contents[i] = &val
		case reflect.Int64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int:
			var val *int64
			contents[i] = &val

		default:
			switch scanType := colType.ScanType(); scanType {
			case reflect.TypeOf(sql.NullInt64{}):
				var val *int64
				contents[i] = &val
			case reflect.TypeOf(sql.NullFloat64{}):
				var val *float64
				contents[i] = &val
				//TODO null mysql
			// case reflect.TypeOf(mysql.NullTime{}):
			// 	var val *time.Time
			// 	contents[i] = &val
			case reflect.TypeOf(sql.RawBytes{}):
				switch colType.DatabaseTypeName() {
				case `CHAR`:
					fallthrough
				case `VARCHAR`:
					fallthrough
				case `TEXT`:
					fallthrough
				case `MEDIUMTEXT`:
					fallthrough
				case `NVARCHAR`:
					var val *string
					contents[i] = &val
				case `DECIMAL`:
					var val *float64
					contents[i] = &val
				case `INT`:
					var val *int64
					contents[i] = &val
				case `JSON`:
					var val *string
					contents[i] = &val
				default:
					var val []byte
					contents[i] = &val
				}
			default:
				switch scanType.Kind() {
				case reflect.Uint64:
					fallthrough
				case reflect.Uint32:
					fallthrough
				case reflect.Uint16:
					fallthrough
				case reflect.Uint8:
					var val *uint64
					contents[i] = &val
				case reflect.Int64:
					fallthrough
				case reflect.Int32:
					fallthrough
				case reflect.Int16:
					fallthrough
				case reflect.Int8:
					var val *int64
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

	}
	return contents
}

//Tx create new Tx when parameter tx is nil and the new Tx will have autocommit enabled. If parameter tx is not null then return tx from parameter
func (c *Connection) Tx(tx *Tx) *Tx {
	if tx == nil {
		return &Tx{db: c.db}
	}
	return tx
}
