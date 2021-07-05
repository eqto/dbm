package db

const (
	noError = iota
	errDuplicate
	errNotFound
	errOther
)

const (
	strNotFound = `no records found`
)

type sqlError struct {
	driver driver
	kind   int
	msg    string
}

func (e *sqlError) Error() string {
	return e.msg
}

func newSQLError(driver driver, kind int) *sqlError {
	s := &sqlError{driver: driver, kind: kind}
	switch kind {
	case errNotFound:
		s.msg = strNotFound
	default:
		s.msg = `Error`
	}
	return s
}

func ErrorDuplicate(e error) bool {
	return isError(e, errDuplicate)
}

//ErrorNotFound to check if no result founds using GetStruct or SelectStruct
func ErrorNotFound(e error) bool {
	return isError(e, errNotFound)
}

func isError(e error, kind int) bool {
	if e == nil {
		return false
	}
	if e, ok := e.(*sqlError); ok {
		if e.kind == 0 {
			if e.driver.isDuplicate(e.msg) {
				e.kind = errDuplicate
			} else {
				e.kind = errOther
			}
		}
		return e.kind == kind
	}
	return false
}

func wrapErr(cn *Connection, e error) *sqlError {
	if e == nil {
		return nil
	}
	return &sqlError{driver: cn.driver, msg: e.Error()}
}
