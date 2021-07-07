package db

import (
	"database/sql"

	"github.com/eqto/go-db/driver"
)

func ExecSelect(bc driver.BuildContents, tx *sql.Tx, query string, params ...interface{}) ([]Resultset, error) {
	rows, e := tx.Query(query, params...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()

	cols, e := rows.Columns()
	if e != nil {
		return nil, e
	}
	colTypes, e := rows.ColumnTypes()
	if e != nil {
		return nil, e
	}

	var results []Resultset

	for rows.Next() {
		contents, e := bc(colTypes)
		if e != nil {
			return nil, e
		}
		if e = rows.Scan(contents...); e != nil {
			rows.Close()
			return nil, e
		}

		rs := Resultset{}
		for key, val := range cols {
			rs[val] = contents[key]
		}
		results = append(results, rs)
	}
	rows.Close()

	return results, nil
}

func ExecGet(bc driver.BuildContents, tx *sql.Tx, query string, params ...interface{}) (Resultset, error) {
	rs, e := ExecSelect(bc, tx, query, params...)
	if e != nil {
		return nil, e
	} else if rs == nil {
		return nil, nil
	}
	return rs[0], nil
}

func ExecInsert(iq driver.InsertQuery, tx *sql.Tx, tableName string, dataMap map[string]interface{}) (*Result, error) {
	length := len(dataMap)
	fields := make([]string, length)
	values := make([]interface{}, length)
	idx := 0
	for name, value := range dataMap {
		fields[idx] = name
		values[idx] = value
		idx++
	}
	query := iq(tableName, fields)
	return Exec(tx, query, values...)
}

func Exec(tx *sql.Tx, query string, params ...interface{}) (*Result, error) {
	res, e := tx.Exec(query, params...)
	if e != nil {
		return nil, e
	}
	return &Result{result: res}, nil
}
