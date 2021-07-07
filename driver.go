package db

import "database/sql"

var (
	drivers = make(map[string]driver)
)

type driver struct {
	name             string
	connectionString func(string, int, string, string, string) string
	isDuplicate      func(string) bool
	buildContents    func([]*sql.ColumnType) ([]interface{}, error)
	insertQuery      func(string, []string) string
}

func Register(name string, rawDriver interface{}) {
	drv := driver{name: name}
	if d, ok := rawDriver.(interface {
		ConnectionString(string, int, string, string, string) string
	}); ok {
		drv.connectionString = d.ConnectionString
	} else {
		return
	}
	if d, ok := rawDriver.(interface {
		IsDuplicate(string) bool
	}); ok {
		drv.isDuplicate = d.IsDuplicate
	} else {
		return
	}
	if d, ok := rawDriver.(interface {
		BuildContents([]*sql.ColumnType) ([]interface{}, error)
	}); ok {
		drv.buildContents = d.BuildContents
	} else {
		return
	}
	if d, ok := rawDriver.(interface {
		InsertQuery(string, []string) string
	}); ok {
		drv.insertQuery = d.InsertQuery
	} else {
		return
	}

	drivers[name] = drv
}
