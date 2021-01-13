package db

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
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

func (s *sqlserverDriver) RegexDuplicate() *regexp.Regexp {
	return regexp.MustCompile(`^mssql: Cannot insert duplicate key.*`)
}

func (s *sqlserverDriver) insertReturnID(tx *Tx, tableName string, fields []string, values []interface{}) (int, error) {
	query := s.insertQuery(tableName, fields) + `; SELECT ID = convert(bigint, SCOPE_IDENTITY())`
	rs, e := tx.Get(query, values...)
	if e != nil {
		return 0, e
	}
	return rs.Int(`ID`), nil
}
