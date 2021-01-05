package db

import (
	"fmt"
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
	return fmt.Sprintf("INSERT INTO `%s`(%s) VALUES(%s)",
		tableName,
		strings.Join(fields, `, `),
		strings.Join(values, `, `))
}

func (m *mysqlDriver) RegexDuplicate() *regexp.Regexp {
	return regexp.MustCompile(`^Duplicate entry.*`)
}
