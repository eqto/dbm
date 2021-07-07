package driver

import (
	"database/sql"
	"fmt"
)

type BuildContents func([]*sql.ColumnType) ([]interface{}, error)
type InsertQuery func(string, []string) string

var (
	drivers = make(map[string]*Driver)
)

type Driver struct {
	Name             string
	ConnectionString func(string, int, string, string, string) string
	IsDuplicate      func(string) bool
	BuildContents    BuildContents
	InsertQuery      InsertQuery
}

func Register(name string, rawDriver interface{}) {
	drv := &Driver{Name: name}
	if d, ok := rawDriver.(interface {
		ConnectionString(string, int, string, string, string) string
	}); ok {
		drv.ConnectionString = d.ConnectionString
	} else {
		return
	}
	if d, ok := rawDriver.(interface {
		IsDuplicate(string) bool
	}); ok {
		drv.IsDuplicate = d.IsDuplicate
	} else {
		return
	}
	if d, ok := rawDriver.(interface {
		BuildContents([]*sql.ColumnType) ([]interface{}, error)
	}); ok {
		drv.BuildContents = d.BuildContents
	} else {
		return
	}
	if d, ok := rawDriver.(interface {
		InsertQuery(string, []string) string
	}); ok {
		drv.InsertQuery = d.InsertQuery
	} else {
		return
	}

	drivers[name] = drv
}

func Get(name string) (*Driver, error) {
	if d, ok := drivers[name]; ok {
		return d, nil
	}
	return nil, fmt.Errorf(`driver '%s' not supported or not registered. Import from github.com/go-db/driver`, name)
}
