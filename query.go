package db

import (
	"fmt"
	"strings"
)

const (
	modeSelect = iota
	modeInsert
	modeDelete
)

type Q struct {
	driverName string
	mode       int
	table      string
	keys       []string
	outputs    []string //sqlserver
	values     []interface{}
}

func (q *Q) InsertInto(table string) *Q {
	q.mode = modeInsert
	q.table = table
	return q
}

//Output used by sqlserver
func (q *Q) Output(keys ...string) *Q {
	if len(keys) == 1 {
		split := strings.Split(keys[0], `,`)
		for idx, str := range split {
			split[idx] = strings.TrimSpace(str)
		}
		q.outputs = split
	} else if len(keys) > 1 {
		q.outputs = keys
	}
	return q
}

func (q *Q) ValueMap(values map[string]interface{}) *Q {
	for key, val := range values {
		q.keys = append(q.keys, key)
		q.values = append(q.values, val)
	}
	return q
}

func (q *Q) Value(key string, value interface{}) *Q {
	q.keys = append(q.keys, key)
	q.values = append(q.values, value)
	return q
}

func (q *Q) String() string {
	s := strings.Builder{}
	switch q.mode {
	case modeInsert:
		values := []string{}
		s.WriteString(`INSERT INTO ` + q.table)
		if len(q.keys) > 0 {
			s.WriteString(`(` + strings.Join(q.keys, `, `) + `)`)
			if q.driverName == `sqlserver` {
				for idx := range q.values {
					values = append(values, fmt.Sprintf(`@p%d`, idx+1))
				}
			} else {
				values = append(values, `?`)
			}
		}
		if len(q.outputs) > 0 {
			for _, output := range q.outputs {
				ucase := strings.ToUpper(output)
				outputs := []string{}
				if strings.HasPrefix(ucase, `INSERTED.`) {
					outputs = append(outputs, `INSERTED`+output[8:])
				} else if strings.HasPrefix(ucase, `DELETED.`) {
					outputs = append(outputs, `DELETED`+output[7:])
				} else {
					switch q.mode {
					case modeInsert:
						outputs = append(outputs, `INSERTED.`+output)
					case modeDelete:
						outputs = append(outputs, `DELETED.`+output)
					}
				}
				s.WriteString(fmt.Sprintf(` OUTPUT %s`, strings.Join(outputs, `, `)))
			}
		}
		s.WriteString(fmt.Sprintf(` VALUES(%s)`, strings.Join(values, `, `)))
	}
	return s.String()
}

func Query(driverName string) *Q {
	return &Q{driverName: driverName}
}
