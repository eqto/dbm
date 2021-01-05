package db

import (
	"regexp"
)

const (
	NoErr = iota
	ErrDuplicate
	ErrOther
)

var (
	duplicateMysql = regexp.MustCompile(`^Duplicate entry.*`)
	duplicateMsSQL = regexp.MustCompile(`^mssql: Cannot insert duplicate key.*`)
)

//SQLError ..
type SQLError interface {
	error
	Kind() int
}

type sqlError struct {
	SQLError
	driver string
	msg    string
}

func (e *sqlError) Error() string {
	return e.msg
}

func (e *sqlError) Kind() int {
	switch e.driver {
	case DriverMySQL:
		if duplicateMysql.MatchString(e.msg) {
			return ErrDuplicate
		}
	case DriverSQLServer:
		if duplicateMsSQL.MatchString(e.msg) {
			return ErrDuplicate
		}
	}
	return ErrOther
}
