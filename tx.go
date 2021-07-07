package db

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

//Tx ...
type Tx struct {
	cn     *Connection
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
		rows, e = t.cn.db.Query(query, params...)
	} else {
		rows, e = t.tx.Query(query, params...)
	}
	if e != nil {
		return nil, wrapErr(t.cn.driver, e)
	}
	defer rows.Close()

	cols, e := rows.Columns()
	if e != nil {
		return nil, wrapErr(t.cn.driver, e)
	}
	colTypes, e := rows.ColumnTypes()
	if e != nil {
		return nil, wrapErr(t.cn.driver, e)
	}

	var results []Resultset

	for rows.Next() {
		contents, e := t.cn.driver.buildContents(colTypes)
		if e != nil {
			return nil, e
		}
		if e = rows.Scan(contents...); e != nil {
			rows.Close()
			return nil, wrapErr(t.cn.driver, e)
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
		return wrapErr(t.cn.driver, e)
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
	if e != nil {
		return nil, wrapErr(t.cn.driver, e)
	} else if rs == nil {
		return nil, nil
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
		return newSQLError(t.cn.driver, errNotFound)
	}

	typeOf = typeOf.Elem()
	return t.assignStruct(dest, t.createFieldMap(typeOf), rs, typeOf)
}

//Exec ...
func (t *Tx) Exec(query string, params ...interface{}) (*Result, error) {
	return exec(t.cn.driver, t.tx, query, params...)
}

func exec(drv *driver, tx *sql.Tx, query string, params ...interface{}) (*Result, error) {
	res, e := tx.Exec(query, params...)
	if e != nil {
		return nil, wrapErr(drv, e)
	}
	return &Result{result: res}, nil
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
	return insert(t.cn.driver, t.tx, tableName, dataMap)
}

func insert(drv *driver, tx *sql.Tx, tableName string, dataMap map[string]interface{}) (*Result, error) {
	length := len(dataMap)
	fields := make([]string, length)
	values := make([]interface{}, length)
	idx := 0
	for name, value := range dataMap {
		fields[idx] = name
		values[idx] = value
		idx++
	}
	query := drv.insertQuery(tableName, fields)
	return exec(drv, tx, query, values...)
}
