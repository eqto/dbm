package dbm

import (
	"database/sql"
)

type Rows struct {
	sqlRows     *sql.Rows
	columnTypes []*sql.ColumnType
	contents    []interface{}
	cols        []string
}

func (r *Rows) Next() bool {
	return r.sqlRows.Next()
}

func (r *Rows) Resultset() (*Resultset, error) {
	if e := r.sqlRows.Scan(r.contents...); e != nil {
		r.sqlRows.Close()
		return nil, e
	}
	rs := Resultset{}
	for key, val := range r.cols {
		rs[val] = r.contents[key]
	}
	return &rs, nil

}

func newRows(driver Driver, sqlRows *sql.Rows) (*Rows, error) {
	rows := &Rows{
		sqlRows: sqlRows,
	}
	types, e := sqlRows.ColumnTypes()
	if e != nil {
		sqlRows.Close()
		return nil, wrapErr(driver, e)
	}
	rows.columnTypes = types
	contents, e := driver.BuildContents(types)
	if e != nil {
		sqlRows.Close()
		return nil, wrapErr(driver, e)
	}
	rows.contents = contents
	cols, e := sqlRows.Columns()
	if e != nil {
		sqlRows.Close()
		return nil, wrapErr(driver, e)
	}
	rows.cols = cols
	return rows, nil
}
