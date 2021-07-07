package db

import "database/sql"

func execSelect(drv *driver, tx *sql.Tx, query string, params ...interface{}) ([]Resultset, error) {
	var rows *sql.Rows
	var e error
	if e != nil {
		return nil, wrapErr(drv, e)
	}
	defer rows.Close()

	cols, e := rows.Columns()
	if e != nil {
		return nil, wrapErr(drv, e)
	}
	colTypes, e := rows.ColumnTypes()
	if e != nil {
		return nil, wrapErr(drv, e)
	}

	var results []Resultset

	for rows.Next() {
		contents, e := drv.buildContents(colTypes)
		if e != nil {
			return nil, e
		}
		if e = rows.Scan(contents...); e != nil {
			rows.Close()
			return nil, wrapErr(drv, e)
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

func execInsert(drv *driver, tx *sql.Tx, tableName string, dataMap map[string]interface{}) (*Result, error) {
	length := len(dataMap)
	fields := make([]string, length)
	values := make([]interface{}, length)
	idx := 0
	for name, value := range dataMap {
		fields[idx] = name
		values[idx] = value
		idx++
	}
	query := drv.insertQuery(tableName, fields)
	return exec(drv, tx, query, values...)
}

func exec(drv *driver, tx *sql.Tx, query string, params ...interface{}) (*Result, error) {
	res, e := tx.Exec(query, params...)
	if e != nil {
		return nil, wrapErr(drv, e)
	}
	return &Result{result: res}, nil
}
