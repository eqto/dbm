package dbqoo

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

//Connection ...
type Connection struct {
	db *sql.DB

	hostname string
	port     int
	username string
	password string
	name     string
	driver   Driver
}

//Connect ...
func (c *Connection) Connect() error {
	db, e := sql.Open(c.driver.Name(), c.driver.DataSourceName(c.hostname, c.port, c.username, c.password, c.name))
	if e != nil {
		return e
	}
	if e := db.Ping(); e != nil {
		return e
	}

	c.db = db
	return nil
}

//Ping ...
func (c *Connection) Ping() error {
	return c.db.Ping()
}

//SetConnMaxLifetime ...
func (c *Connection) SetConnMaxLifetime(duration time.Duration) {
	c.db.SetConnMaxLifetime(duration)
}

//SetMaxIdleConns ...
func (c *Connection) SetMaxIdleConns(max int) {
	c.db.SetMaxIdleConns(max)
}

//SetMaxOpenConns ...
func (c *Connection) SetMaxOpenConns(max int) {
	c.db.SetMaxOpenConns(max)
}

//Begin ...
func (c *Connection) Begin() (*Tx, error) {
	sqlTx, e := c.db.Begin()
	if e != nil {
		return nil, e
	}
	return &Tx{drv: c.driver, sqlTx: sqlTx}, nil
}

//MustBegin ...
func (c *Connection) MustBegin() *Tx {
	tx, e := c.Begin()
	if e != nil {
		return nil
	}
	return tx
}

//MustExec ...
func (c *Connection) MustExec(query string, params ...interface{}) *Result {
	result, e := c.Exec(query, params...)
	if e != nil {
		panic(e)
	}
	return result
}

//Exec ...
func (c *Connection) Exec(query string, params ...interface{}) (*Result, error) {
	tx, e := c.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Recover()
	return tx.Exec(query, params...)
}

//Get ...
func (c *Connection) Get(query string, params ...interface{}) (Resultset, error) {
	tx, e := c.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Recover()

	return tx.Get(query, params...)
}

//MustGet ...
func (c *Connection) MustGet(query string, params ...interface{}) Resultset {
	rs, e := c.Get(query, params...)
	if e != nil {
		panic(e)
	}
	return rs
}

//GetStruct ...
func (c *Connection) GetStruct(dest interface{}, query string, params ...interface{}) error {
	tx, e := c.Begin()
	if e != nil {
		return e
	}
	defer tx.Recover()

	return tx.GetStruct(dest, query, params...)
}

//MustSelect ...
func (c *Connection) MustSelect(query string, params ...interface{}) []Resultset {
	rs, e := c.Select(query, params...)
	if e != nil {
		panic(e)
	}
	return rs
}

//Select ...
func (c *Connection) Select(query string, params ...interface{}) ([]Resultset, error) {
	tx, e := c.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Recover()
	return tx.Select(query, params...)
}

//SelectStruct ...
func (c *Connection) SelectStruct(dest interface{}, query string, params ...interface{}) error {
	tx, e := c.Begin()
	if e != nil {
		return e
	}
	defer tx.Recover()
	return tx.SelectStruct(dest, query, params...)
}

//MustInsert ...
func (c *Connection) MustInsert(tableName string, dataMap map[string]interface{}) *Result {
	result, e := c.Insert(tableName, dataMap)
	if e != nil {
		panic(e)
	}
	return result
}

//Insert ...
func (c *Connection) Insert(tableName string, dataMap map[string]interface{}) (*Result, error) {
	tx, e := c.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Recover()
	return tx.Insert(tableName, dataMap)
}

//EnumValues return enum values, parameter field using dot notation. Ex: profile.gender , returning ['male', 'female']
func (c *Connection) EnumValues(field string) ([]string, error) {
	cols := strings.Split(field, `.`)
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

//Close ...
func (c *Connection) Close() error {
	if c.db == nil {
		return nil
	}
	return c.db.Close()
}

func newConnection(driverName, hostname string, port int, username, password, name string) (*Connection, error) {
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf(`invalid port %d`, port)
	}

	drv, e := getDriver(driverName)
	if e != nil {
		return nil, e
	}
	return &Connection{
		driver:   drv,
		hostname: hostname, port: port,
		username: username, password: password,
		name: name,
	}, nil

}
