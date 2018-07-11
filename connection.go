/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2017-12-17 22:16:39
*/

package db

import (
	"bytes"
	"regexp"
	"strings"
	"database/sql"
)

//Connection ...
type Connection struct	{
	db			*sql.DB
}

//Begin ...
func (c *Connection) Begin() (*Tx, error)	{
	tx, e := c.db.Begin()
	if e != nil	{
		return nil, e
	}
	return &Tx{ tx: tx }, nil
}
	
//MustBegin ...
func (c *Connection) MustBegin() *Tx	{
	tx, e := c.Begin()
	if e != nil	{
		return nil
	}
	return tx
}

//MustExec ...
func (c *Connection) MustExec(query string, params...interface{}) *Result	{
	result, e := c.Exec(query, params...)
	if e != nil	{
		panic(e)
	}
	return result
}
//Exec ...
func (c *Connection) Exec(query string, params...interface{}) (*Result, error)	{
	tx, e := c.Begin()
	if e != nil	{
		return nil, e
	}
	defer func ()  {
		if r := recover(); r != nil	{
			tx.Rollback()
			panic(r)
		} else {
			tx.Commit()
		}
	}()
	return tx.Exec(query, params...)
}

//MustGet ...
func (c *Connection) MustGet(query string, params...interface{}) Resultset	{
	rs, e := c.Get(query, params...)
	if e != nil	{
		panic(e)
	}
	return rs
}
	
//Get ...
func (c *Connection) Get(query string, params...interface{}) (Resultset, error)	{
	tx, e := c.Begin()
	if e != nil	{
		return nil, e
	}
	defer func ()  {
		if r := recover(); r != nil	{
			tx.Rollback()
		}
		tx.Commit()
	}()
	return tx.Get(query, params...)
}	

//GetStruct ...
func (c *Connection) GetStruct(dest interface{}, query string, params...interface{}) error	{
	tx, e := c.Begin()
	if e != nil	{
		return e
	}
	defer func ()  {
		if r := recover(); r != nil	{
			tx.Rollback()
		}
		tx.Commit()
	}()

	return tx.GetStruct(dest, query, params...)
}
	
//MustSelect ...
func (c *Connection) MustSelect(query string, params...interface{}) []Resultset	{
	rs, e := c.Select(query, params...)
	if e != nil	{
		panic(e)
	}
	return rs
}

//Select ...
func (c *Connection) Select(query string, params...interface{}) ([]Resultset, error)	{
	tx, e := c.Begin()
	if e != nil	{
		return nil, e
	}
	defer func ()  {
		if r := recover(); r != nil	{
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	return tx.Select(query, params...)
}

//SelectStruct ...
func (c *Connection) SelectStruct(dest interface{}, query string, params...interface{}) error	{
	tx, e := c.Begin()
	if e != nil	{
		return e
	}
	defer func ()  {
		if r := recover(); r != nil	{
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	return tx.SelectStruct(dest, query, params...)
}

//MustInsert ...
func (c *Connection) MustInsert(tableName string, dataMap map[string]interface{}) *Result	{
	result, e := c.Insert(tableName, dataMap)
	if e != nil	{
		panic(e)
	}
	return result
}

//Insert ...
func (c *Connection) Insert(tableName string, dataMap map[string]interface{}) (*Result, error)	{
	var names []string
	var questionMarks []string
	var values []interface{}

	for name, value := range dataMap	{
		names = append(names, name)
		values = append(values, value)
		questionMarks = append(questionMarks, `?`)
	}
	var buffer bytes.Buffer
	buffer.WriteString(`INSERT INTO `)
	buffer.WriteString(tableName)
	buffer.WriteString(`(` + strings.Join(names, `, `) + `)`)
	buffer.WriteString(` VALUES(` + strings.Join(questionMarks, `, `) + `)`)

	return c.Exec(buffer.String(), values...)
}

//GetEnumValues ...
func (c *Connection) GetEnumValues(field string) ([]string, error)	{
	cols := strings.Split(field, `.`)

	enum, e := c.Get(`SELECT column_type FROM information_schema.columns WHERE table_name = ? 
		AND column_name = ?`, cols[0], cols[1])
	if e != nil	{
		return nil, e
	}
	regexEnum := regexp.MustCompile(`'[a-zA-Z0-9]+'`)

	values := regexEnum.FindAllString(*enum.GetString(`column_type`), -1)

	for i := 0; i < len(values); i++	{
		values[i] = strings.Trim(values[i], `'`)
	}
	return values, nil
}

//Close ...
func (c *Connection) Close() error	{
	if c.db == nil	{
		return nil
	}
	return c.db.Close()
}