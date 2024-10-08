package dbm

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// Connection ...
type Connection struct {
	db *sql.DB

	cfg Config
	drv Driver
}

func nilCheck(c *Connection) error {
	if c.db == nil {
		return ErrConnectionNotFound
	}
	if c.drv == nil {
		return ErrDriverNotDefined
	}
	return nil
}

func (c *Connection) applyOptions(opts ...Options) {
	for _, opt := range opts {
		opt(c)
	}
}

// Connect ...
func (c *Connection) Connect(opts ...Options) error {
	if c.drv == nil {
		return ErrDriverNotDefined
	}
	c.applyOptions(opts...)
	db, e := sql.Open(c.cfg.DriverName, c.drv.DataSourceName(c.cfg))
	if e != nil {
		return e
	}
	c.db = db

	c.applyOptions([]Options{
		OptionMaxIdleTime(60 * time.Second),
		OptionMaxLifetime(5 * time.Minute),
		OptionMaxIdle(2)}...)

	c.applyOptions(opts...)

	if e := db.Ping(); e != nil {
		return e
	}

	return nil
}

func (c *Connection) Close() {
	if c.db == nil {
		c.db.Close()
		c.db = nil
	}
}

func (c *Connection) Driver() Driver {
	return c.drv
}

// Ping ...
func (c *Connection) Ping() error {
	if e := nilCheck(c); e != nil {
		return e
	}
	return c.db.Ping()
}

// SetConnMaxLifetime ...
func (c *Connection) SetConnMaxLifetime(duration time.Duration) {
	c.db.SetConnMaxLifetime(duration)
}

// SetMaxIdleConns ...
func (c *Connection) SetMaxIdleConns(max int) {
	c.db.SetMaxIdleConns(max)
}

// SetMaxOpenConns ...
func (c *Connection) SetMaxOpenConns(max int) {
	c.db.SetMaxOpenConns(max)
}

// Begin ...
func (c *Connection) Begin() (*Tx, error) {
	if e := nilCheck(c); e != nil {
		return nil, e
	}
	sqlTx, e := c.db.Begin()
	if e != nil {
		return nil, e
	}
	return &Tx{driver: c.drv, sqlTx: sqlTx}, nil
}

// MustBegin ...
func (c *Connection) MustBegin() *Tx {
	tx, e := c.Begin()
	if e != nil {
		return nil
	}
	return tx
}

// MustExec ...
func (c *Connection) MustExec(query string, args ...interface{}) *Result {
	result, e := c.Exec(query, args...)
	if e != nil {
		panic(e)
	}
	return result
}

// Exec ...
func (c *Connection) Exec(query string, args ...interface{}) (*Result, error) {
	if e := nilCheck(c); e != nil {
		return nil, e
	}
	return exec(c.drv, c.db.Exec, query, args...)
}

// Get ...
func (c *Connection) Get(query string, args ...interface{}) (Resultset, error) {
	rs, e := c.Select(query, args...)
	if e != nil {
		return nil, e
	} else if rs == nil {
		return nil, nil
	}
	return rs[0], nil
}

// MustGet ...
func (c *Connection) MustGet(query string, args ...interface{}) Resultset {
	rs, e := c.Get(query, args...)
	if e != nil {
		panic(e)
	}
	return rs
}

// GetStruct ...
func (c *Connection) GetStruct(dest interface{}, query string, args ...interface{}) error {
	if dest == nil {
		return errors.New(`parameter dest is nil`)
	}
	typeOf := reflect.TypeOf(dest)
	if typeOf.Kind() != reflect.Ptr {
		return errors.New(`parameter dest is not a pointer`)
	}

	rs, e := c.Get(query, args...)
	if e != nil {
		return e
	} else if rs == nil {
		return errors.New(errRecordNotFound)
	}

	typeOf = typeOf.Elem()
	return assignStruct(dest, createFieldMap(typeOf), rs, typeOf)
}

// MustSelect ...
func (c *Connection) MustSelect(query string, args ...interface{}) []Resultset {
	rs, e := c.Select(query, args...)
	if e != nil {
		panic(e)
	}
	return rs
}

func (c *Connection) Query(query string, args ...interface{}) (*Rows, error) {
	if e := nilCheck(c); e != nil {
		return nil, e
	}
	rows, e := c.db.Query(query, args...)
	if e != nil {
		return nil, e
	}
	return newRows(c.drv, rows)
}

// Select ...
func (c *Connection) Select(query string, args ...interface{}) ([]Resultset, error) {
	if e := nilCheck(c); e != nil {
		return nil, e
	}
	return execQuery(c.drv, c.db.Query, query, args...)
}

// SelectStruct ...
func (c *Connection) SelectStruct(dest interface{}, query string, args ...interface{}) error {
	if e := nilCheck(c); e != nil {
		return e
	}
	return execQueryStruct(c.drv, c.Select, dest, query, args...)
}

// MustInsert ...
func (c *Connection) MustInsert(tableName string, dataMap map[string]interface{}) *Result {
	result, e := c.Insert(tableName, dataMap)
	if e != nil {
		panic(e)
	}
	return result
}

// Insert ...
func (c *Connection) Insert(tableName string, dataMap map[string]interface{}) (*Result, error) {
	if e := nilCheck(c); e != nil {
		return nil, e
	}
	length := len(dataMap)
	fields := make([]string, length)
	values := []interface{}{}
	placeholders := []string{}
	idx := 0
	for name, value := range dataMap {
		fields[idx] = name
		if val, ok := value.(SQLStatement); ok {
			placeholders = append(placeholders, val.statement)
		} else {
			placeholders = append(placeholders, `?`)

			values = append(values, value)
		}
		idx++
	}
	q := InsertInto(tableName, strings.Join(fields, `, `)).Values(strings.Join(placeholders, `, `))
	return c.Exec(c.drv.StatementString(q), values...)
}

func (c *Connection) SQL(stmt interface{}) string {
	if c.drv != nil {
		return c.drv.StatementString(stmt)
	}
	return ``
}

// EnumValues return enum values, parameter field using dot notation. Ex: profile.gender , returning ['male', 'female']
func (c *Connection) EnumValues(field string) ([]string, error) {
	cols := strings.Split(field, `.`)
	if len(cols) != 2 {
		return nil, errors.New(`invalid parameter field:` + field)
	}
	enum, e := c.Get(`SELECT column_type FROM information_schema.columns WHERE table_name = ?
		AND column_name = ?`, cols[0], cols[1])
	if e != nil {
		return nil, e
	}
	regexEnum := regexp.MustCompile(`'[a-zA-Z0-9]+'`)

	values := regexEnum.FindAllString(enum.String(`column_type`), -1)

	for i := 0; i < len(values); i++ {
		values[i] = strings.Trim(values[i], `'`)
	}
	return values, nil
}

// Close ...
func (c *Connection) Close() error {
	if c.db == nil {
		return nil
	}
	return c.db.Close()
}
