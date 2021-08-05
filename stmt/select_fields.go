package stmt

import "strings"

type SelectFields struct {
	value Fields
}

func (f SelectFields) From(table string) *Select {
	split := strings.SplitN(table, ` `, 2)
	tables := Tables{}
	if len(split) == 2 {
		tables.add(strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), 0, ``)
	} else {
		tables.add(strings.TrimSpace(table), ``, 0, ``)
	}
	return &Select{fields: f, tables: tables}
}
