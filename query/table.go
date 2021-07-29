/*
 * @Author: tuxer
 * @Date: 2021-07-29 12:42:08
 * @Last Modified by: tuxer
 * @Last Modified time: 2021-07-29 13:51:59
 */

package query

import (
	"strings"
)

type Table struct {
	Name      string
	Alias     string
	Condition string

	Fields []Field
}

func parseTable(query string, fields []Field) Table {
	t := Table{}
	split := strings.SplitN(query, ` `, 2)
	prefix := ``
	if len(split) == 2 {
		t.Name, t.Alias = strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
		prefix = t.Alias + `.`
	} else {
		t.Name = strings.TrimSpace(query)
		prefix = t.Name + `.`
	}
	for _, f := range fields {
		if strings.HasPrefix(f.Name, prefix) {
			t.Fields = append(t.Fields, f)
		}
	}
	return t
}
