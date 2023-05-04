package dbm

const (
	errRecordNotFound     = `record not found`
	ErrConnectionNotFound = Error(`connection not found`)
	ErrDriverNotDefined   = Error(`database driver not defined`)
)

type Error string

func (e Error) Error() string {
	return string(e)
}

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

func IsErrNotFound(e error) bool {
	return e.Error() == errRecordNotFound
}

func wrapErr(drv Driver, e error) *sqlError {
	if e == nil {
		return nil
	}
	return &sqlError{drv: drv, e: e}
}
