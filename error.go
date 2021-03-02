package db

const (
	NoErr = iota
	ErrNotFound
	ErrDuplicate
	ErrOther
)

const (
	strNotFound = `no records found`
)

//SQLError ..
type SQLError interface {
	error
	Kind() int
}

type sqlError struct {
	SQLError
	driver driver
	kind   int
	msg    string
}

func (e *sqlError) Error() string {
	return e.msg
}

func (e *sqlError) Kind() int {
	if e.kind > 0 {
		return e.kind
	}
	if e.driver.regexDuplicate().MatchString(e.msg) {
		return ErrDuplicate
	}
	return ErrOther
}

func newSQLError(driver driver, kind int) SQLError {
	s := &sqlError{driver: driver, kind: kind}
	switch kind {
	case ErrNotFound:
		s.msg = strNotFound
	default:
		s.msg = `Error`
	}
	return s
}
