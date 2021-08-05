package dbm

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
	drv  Driver
	kind int
	msg  string
}

func (e *sqlError) Error() string {
	return e.msg
}

func newSQLError(drv Driver, kind int) *sqlError {
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
			if e.drv.IsDuplicate(e.msg) {
				e.kind = errDuplicate
			} else {
				e.kind = errOther
			}
		}
		return e.kind == kind
	}
	return false
}

func wrapErr(drv Driver, e error) *sqlError {
	if e == nil {
		return nil
	}
	return &sqlError{drv: drv, msg: e.Error()}
}
