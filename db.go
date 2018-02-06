/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2017-12-16 16:36:19
*/

package db

import (
	"regexp"
	"database/sql"
)

var (
	defaultCn			*Connection
	regexStringColType	*regexp.Regexp
)

//NewConnection ...
func NewConnection(host string, username string, password string, databaseName string) (*Connection, error)	{
	db, e := sql.Open(`mysql`, username + `:` + password + `@tcp(` + host + `)/` + 
		databaseName + `?parseTime=true&loc=Local&collation=latin1_general_ci`)

	if e != nil	{
		return nil, e
	}
	if e := db.Ping(); e != nil	{
		return nil, e
	}

	defaultCn = &Connection{db: db}
	return defaultCn, nil
}

func getRegex() *regexp.Regexp	{
	if regexStringColType == nil	{
		regexStringColType, _ = regexp.Compile(`(?i)^.*CHAR|TEXT$`)
	}
	return regexStringColType
}
