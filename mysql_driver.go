package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type mysqlDriver struct {
	driver
	params
}

func (m *mysqlDriver) connectionString() string {
	return fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`,
		m.username, m.password,
		m.hostname, m.port,
		m.name,
	)
}

func (m *mysqlDriver) kind() string {
	return `mysql`
}

func (m *mysqlDriver) insertQuery(tableName string, fields []string) string {
	values := make([]string, len(fields))
	for i := range values {
		values[i] = `?`
	}
	return fmt.Sprintf("INSERT INTO `%s`(`%s`) VALUES(%s)",
		tableName,
		strings.Join(fields, "`, `"),
		strings.Join(values, `, `))
}

func (m *mysqlDriver) regexDuplicate() *regexp.Regexp {
	return regexp.MustCompile(`^Duplicate entry.*`)
}

func (m *mysqlDriver) insertReturnID(tx *Tx, tableName string, fields []string, values []interface{}) (int, error) {
	res, e := tx.Exec(tableName, fields, values)
	if e != nil {
		return 0, e
	}
	id, e := res.LastInsertID()
	if e != nil {
		return 0, e
	}
	return id, nil
}

func (m *mysqlDriver) buildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
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
			vals[idx] = new([]byte)
		case reflect.Struct:
			switch scanType.Name() {
			case `NullInt64`:
				vals[idx] = new(sql.NullInt64)
			case `NullFloat64`:
				vals[idx] = new(sql.NullFloat64)
			case `NullTime`:
				vals[idx] = new(sql.NullTime)
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}
