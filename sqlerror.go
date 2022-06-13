package dbm

type sqlError struct {
	drv Driver
	e   error
}

func (e *sqlError) Error() string {
	return e.e.Error()
}

func IsErrDuplicate(e error) bool {
	if e, ok := e.(*sqlError); ok {
		return e.drv.IsDuplicate(e.e)
	}
	return false
}

func wrapErr(drv Driver, e error) *sqlError {
	if e == nil {
		return nil
	}
	return &sqlError{drv: drv, e: e}
}
