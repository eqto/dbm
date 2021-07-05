package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type sqlserverDriver struct {
	driver
	params
}

func (s *sqlserverDriver) connectionString() string {
	u := url.URL{
		Scheme:   s.kind(),
		User:     url.UserPassword(s.username, s.password),
		Host:     fmt.Sprintf("%s:%d", s.hostname, s.port),
		RawQuery: s.name,
	}
	return u.String()
}

func (s *sqlserverDriver) kind() string {
	return `sqlserver`
}

func (s *sqlserverDriver) insertQuery(tableName string, fields []string) string {
	values := make([]string, len(fields))
	for i := range values {
		values[i] = fmt.Sprintf(`@p%d`, i+1)
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		tableName,
		strings.Join(fields, `, `),
		strings.Join(values, `, `))
}

func (s *sqlserverDriver) isDuplicate(msg string) bool {
	return regexp.MustCompile(`.*Cannot insert duplicate key.*`).MatchString(msg)
}

func (s *sqlserverDriver) insertReturnID(tx *Tx, tableName string, fields []string, values []interface{}) (int, error) {
	query := s.insertQuery(tableName, fields) + `; SELECT ID = convert(bigint, SCOPE_IDENTITY())`
	rs, e := tx.Get(query, values...)
	if e != nil {
		return 0, e
	}
	return rs.Int(`ID`), nil
}

func (s *sqlserverDriver) buildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
	vals := make([]interface{}, len(colTypes))
	for idx, colType := range colTypes {
		scanType := colType.ScanType()
		switch scanType.Kind() {
		case reflect.Int64:
			vals[idx] = new(*int64)
		case reflect.Bool:
			vals[idx] = new(*bool)
		case reflect.String:
			vals[idx] = new(*string)
		case reflect.Struct:
			switch scanType.Name() {
			case `Time`:
				vals[idx] = new(*time.Time)
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}
