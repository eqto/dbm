package db

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

//Tx ...
type Tx struct {
	drv    *driver
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
	return execSelect(t.drv, t.tx, query, params...)
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
		return wrapErr(t.drv, e)
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
		return nil, wrapErr(t.drv, e)
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
		return newSQLError(t.drv, errNotFound)
	}

	typeOf = typeOf.Elem()
	return t.assignStruct(dest, t.createFieldMap(typeOf), rs, typeOf)
}

//Exec ...
func (t *Tx) Exec(query string, params ...interface{}) (*Result, error) {
	return exec(t.drv, t.tx, query, params...)
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
	return execInsert(t.drv, t.tx, tableName, dataMap)
}
