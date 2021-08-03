package db

import (
	"database/sql"
	"fmt"
)

var (
	driverMap = make(map[string]Driver)
)

type Driver interface {
	Name() string
	DataSourceName(string, int, string, string, string) string
	Query(stmt interface{}) string
	IsDuplicate(string) bool
	BuildContents([]*sql.ColumnType) ([]interface{}, error)
}

func Register(name string, driver Driver) {
	driverMap[name] = driver
}

func getDriver(name string) (Driver, error) {
	if d, ok := driverMap[name]; ok {
		return d, nil
	}
	return nil, fmt.Errorf(`driver '%s' not supported or not registered. Import from github.com/dbq/driver`, name)
}
