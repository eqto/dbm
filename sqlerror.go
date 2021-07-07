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
	drv  *driver
	kind int
	msg  string
}

func (e *sqlError) Error() string {
	return e.msg
}

func newSQLError(drv *driver, kind int) *sqlError {
	s := &sqlError{drv: drv, kind: kind}
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
			if e.drv.isDuplicate(e.msg) {
				e.kind = errDuplicate
			} else {
				e.kind = errOther
			}
		}
		return e.kind == kind
	}
	return false
}

func wrapErr(drv *driver, e error) *sqlError {
	if e == nil {
		return nil
	}
	return &sqlError{drv: drv, msg: e.Error()}
}
