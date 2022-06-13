package sqlserver

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/eqto/dbm"
	"github.com/eqto/dbm/stmt"
)

func init() {
	dbm.Register(`sqlserver`, &Driver{})
}

type Driver struct {
	dbm.Driver
}

func (Driver) Name() string {
	return `sqlserver`
}

func (Driver) SanitizeParams(values []interface{}) []interface{} {
	for idx := range values {
		switch value := values[idx].(type) {
		case string:
			values[idx] = mssql.VarChar(value)
		}
	}
	return values
}

func (Driver) StatementString(s interface{}) string {
	s = stmt.StatementOf(s)
	switch s := s.(type) {
	case *stmt.Select:
		return selectStatement(s)
	case *stmt.Insert:
		return insertStatement(s)
	case *stmt.Update:
		return updateStatement(s)
	}
	return ``
}

func (Driver) DataSourceName(hostname string, port int, username, password, name string) string {
	u := url.URL{
		Scheme:   `sqlserver`,
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%d", hostname, port),
		RawQuery: fmt.Sprintf(`database=%s&app+name=dbm&TrustServerCertificate=true`, name),
	}
	return u.String()
}
func (Driver) IsDuplicate(msg string) bool {
	return regexp.MustCompile(`.*Cannot insert duplicate key.*`).MatchString(msg)
}

func (Driver) BuildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
	vals := make([]interface{}, len(colTypes))
	for idx, colType := range colTypes {
		scanType := colType.ScanType()
		switch scanType.Kind() {
		case reflect.Int64:
			vals[idx] = new(*int64)
		case reflect.Bool:
			vals[idx] = new(*bool)
		case reflect.String:
			vals[idx] = new(*string)
		case reflect.Struct:
			switch scanType.Name() {
			case `Time`:
				vals[idx] = new(*time.Time)
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}
