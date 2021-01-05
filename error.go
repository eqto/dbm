package db

const (
	NoErr = iota
	ErrDuplicate
	ErrOther
)

//SQLError ..
type SQLError interface {
	error
	Kind() int
}

type sqlError struct {
	SQLError
	driver driver
	msg    string
}

func (e *sqlError) Error() string {
	return e.msg
}

func (e *sqlError) Kind() int {
	if e.driver.RegexDuplicate().MatchString(e.msg) {
		return ErrDuplicate
	}
	return ErrOther
}
