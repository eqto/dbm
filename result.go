package dbm

import (
	"database/sql"
)

// Result ...
type Result struct {
	result sql.Result
}

// LastInsertID ...
func (r *Result) LastInsertID() (ID int, e error) {
	if r.result == nil {
		return 0, nil
	}
	id, e := r.result.LastInsertId()
	if e != nil {
		return 0, e
	}
	return int(id), nil
}

// MustLastInsertID ...
func (r *Result) MustLastInsertID() int {
	id, e := r.LastInsertID()
	if e != nil {
		panic(e)
	}
	return id
}

// RowsAffected ...
func (r *Result) RowsAffected() (int, error) {
	if r.result == nil {
		return 0, nil
	}
	val, e := r.result.RowsAffected()
	if e != nil {
		return 0, e
	}
	return int(val), nil
}

// MustRowsAffected ...
func (r *Result) MustRowsAffected() int {
	row, e := r.RowsAffected()
	if e != nil {
		panic(e)
	}
	return int(row)
}
