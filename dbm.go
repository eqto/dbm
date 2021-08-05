package dbm

import "github.com/eqto/dbm/stmt"

//Connect ...
func Connect(driver, host string, port int, username, password, name string) (*Connection, error) {
	cn, e := newConnection(driver, host, port, username, password, name)
	if e != nil {
		return nil, e
	}
	if e := cn.Connect(); e != nil {
		return nil, e
	}
	return cn, nil
}

func Select(fields string) *stmt.SelectFields {
	return stmt.Build().Select(fields)
}

func InsertInto(table, fields string) *stmt.Insert {
	return stmt.Build().InsertInto(table, fields)
}

func Update(table string) *stmt.Update {
	return stmt.Build().Update(table)
}
