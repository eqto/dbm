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

func (Driver) DataSourceName(cfg dbm.Config) string {
	query := fmt.Sprintf(`database=%s&app+name=dbm&TrustServerCertificate=true`, cfg.Name)
	if cfg.DisableEncryption {
		query += `&encrypt=disable`
	}
	u := url.URL{
		Scheme:   `sqlserver`,
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port),
		RawQuery: query,
	}
	return u.String()
}
func (Driver) IsDuplicate(e error) bool {
	return regexp.MustCompile(`.*Cannot insert duplicate key.*`).MatchString(e.Error())
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
