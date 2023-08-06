package dbm

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

type Tx struct {
	driver Driver
	sqlTx  *sql.Tx
	finish bool
}

func (t *Tx) MustRecover() {
	if r := recover(); r != nil {
		t.Rollback()
		panic(r)
	}
	t.Commit()
}

func (t *Tx) Recover() {
	if r := recover(); r != nil {
		t.Rollback()
	} else {
		t.Commit()
	}
}

func (t *Tx) Commit() error {
	if t == nil {
		return errors.New(`Connection already closed`)
	}
	if t.finish {
		return nil
	}
	t.finish = true
	return t.sqlTx.Commit()
}

func (t *Tx) Rollback() error {
	if t == nil {
		return errors.New(`Connection already closed`)
	}
	if t.finish {
		return nil
	}
	t.finish = true
	return t.sqlTx.Rollback()
}

func (t *Tx) MustInsert(tableName string, dataMap map[string]interface{}) *Result {
	if rs, e := t.Insert(tableName, dataMap); e != nil {
		panic(e)
	} else {
		return rs
	}
}

func (t *Tx) MustSelect(query string, args ...interface{}) []Resultset {
	if rs, e := t.Select(query, args...); e != nil {
		panic(e)
	} else {
		return rs
	}
}

func (t *Tx) Select(query string, args ...interface{}) ([]Resultset, error) {
	return execQuery(t.driver, t.sqlTx.Query, query, args...)
}

func (t *Tx) SelectStruct(dest interface{}, query string, args ...interface{}) error {
	return execQueryStruct(t.driver, t.Select, dest, query, args...)
}

func (t *Tx) Get(query string, args ...interface{}) (Resultset, error) {
	rs, e := t.Select(query, args...)
	if e != nil {
		return nil, e
	} else if rs == nil {
		return nil, nil
	}
	return rs[0], nil
}

func (t *Tx) MustGet(query string, args ...interface{}) Resultset {
	rs, e := t.Get(query, args...)
	if e != nil {
		panic(e)
	}
	return rs
}

func (t *Tx) GetStruct(dest interface{}, query string, args ...interface{}) error {
	typeOf := reflect.TypeOf(dest)
	if typeOf.Kind() != reflect.Ptr {
		return errors.New(`dest is not a pointer`)
	}

	rs, e := t.Get(query, args...)
	if e != nil {
		return e
	} else if rs == nil {
		return errors.New(errRecordNotFound)
	}

	typeOf = typeOf.Elem()
	return assignStruct(dest, createFieldMap(typeOf), rs, typeOf)
}

func (t *Tx) Exec(query string, args ...interface{}) (*Result, error) {
	return exec(t.driver, t.sqlTx.Exec, query, args...)
}

func (t *Tx) MustExec(query string, args ...interface{}) *Result {
	result, e := t.Exec(query, args...)
	if e != nil {
		panic(e)
	}
	return result
}

func (t *Tx) Insert(tableName string, dataMap map[string]interface{}) (*Result, error) {
	length := len(dataMap)
	fields := make([]string, length)
	values := []interface{}{}
	placeholders := []string{}
	idx := 0
	for name, value := range dataMap {
		fields[idx] = name
		if val, ok := value.(SQLStatement); ok {
			placeholders = append(placeholders, val.statement)
		} else {
			placeholders = append(placeholders, `?`)

			values = append(values, value)
		}
		idx++
	}
	q := InsertInto(tableName, strings.Join(fields, `, `)).Values(strings.Join(placeholders, `, `))
	return t.Exec(t.driver.StatementString(q), values...)
}
