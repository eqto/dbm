package db

const (
	NoErr = iota
	ErrNotFound
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
	if e.driver.regexDuplicate().MatchString(e.msg) {
		return ErrDuplicate
	}
	if regex := e.driver.regexRecordNotFound(); regex != nil && regex.MatchString(e.msg) {
		return ErrNotFound
	}
	return ErrOther
}
