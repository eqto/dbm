package db

import (
	"regexp"
)

var (
	regexStringColType *regexp.Regexp
)

type driver interface {
	connectionString() string
	kind() string
	insertQuery(tableName string, fields []string) string
	insertReturnID(tx *Tx, tableName string, fields []string, values []interface{}) (int, error)

	regexDuplicate() *regexp.Regexp
}

func getRegex() *regexp.Regexp {
	if regexStringColType == nil {
		regexStringColType, _ = regexp.Compile(`(?i)^.*CHAR|.*TEXT$`)
	}
	return regexStringColType
}
