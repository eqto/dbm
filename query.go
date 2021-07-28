package db

import (
	"fmt"
	"strings"
)

type Table struct {
	name, alias string
}

func (t Table) String() string {
	return strings.TrimSpace(fmt.Sprintf(`%s %s`, t.name, t.alias))
}

type Field struct {
	table, name, alias string
}

func (f Field) String() string {
	sb := strings.Builder{}
	if f.table != `` {
		sb.WriteString(fmt.Sprintf(`%s.`, f.table))
	}
	sb.WriteString(fmt.Sprintf(`%s`, f.name))
	if f.alias != `` {
		sb.WriteString(fmt.Sprintf(` AS %s`, f.alias))
	}
	return sb.String()
}
