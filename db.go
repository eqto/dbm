/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2017-12-16 16:36:19
 */

package db

import (
	"database/sql"
	"fmt"
	"regexp"

	//mysql driver ...
	_ "github.com/go-sql-driver/mysql"
)

var (
	lastCn             *Connection
	regexStringColType *regexp.Regexp
)

//NewConnection ...
func NewConnection(host string, port int, username, password, databaseName string) (*Connection, error) {
	db, e := sql.Open(`mysql`,
		fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`, username, password, host, port, databaseName))

	if e != nil {
		return nil, e
	}
	if e := db.Ping(); e != nil {
		return nil, e
	}

	lastCn = &Connection{db: db}
	return lastCn, nil
}

func getRegex() *regexp.Regexp {
	if regexStringColType == nil {
		regexStringColType, _ = regexp.Compile(`(?i)^.*CHAR|.*TEXT$`)
	}
	return regexStringColType
}
