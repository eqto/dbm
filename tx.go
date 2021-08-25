package dbm

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

//Tx ...
type Tx struct {
	drv    Driver
	sqlTx  *sql.Tx
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
	if t.finish {
		return nil
	}
	t.finish = true
	return t.sqlTx.Commit()
}

//Rollback ...
func (t *Tx) Rollback() error {
	if t.finish {
		return nil
	}
	t.finish = true
	return t.sqlTx.Rollback()
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
	return execQuery(t.drv, t.sqlTx.Query, query, params...)
}

//SelectStruct ...
func (t *Tx) SelectStruct(dest interface{}, query string, params ...interface{}) error {
	return execQueryStruct(t.Select, dest, query, params...)
}

//Get ...
func (t *Tx) Get(query string, params ...interface{}) (Resultset, error) {
	rs, e := t.Select(query, params...)
	if e != nil {
		return nil, e
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
	return assignStruct(dest, createFieldMap(typeOf), rs, typeOf)
}

//Exec ...
func (t *Tx) Exec(query string, params ...interface{}) (*Result, error) {
	res, e := t.sqlTx.Exec(query, params...)
	if e != nil {
		return nil, wrapErr(t.drv, e)
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
	length := len(dataMap)
	fields := make([]string, length)
	values := make([]interface{}, length)
	idx := 0
	for name, value := range dataMap {
		fields[idx] = name
		values[idx] = value
		idx++
	}
	q := InsertInto(tableName, strings.Join(fields, `, `))
	return t.Exec(t.drv.StatementString(q), values...)
}
