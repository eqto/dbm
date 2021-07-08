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
	driverName   string
	mode         int
	table        string
	fields       []string
	keys         []string
	outputs      []string //sqlserver
	values       []interface{}
	wheres       []string
	start, count int
}

func (q *Q) InsertInto(table string) *Q {
	q.mode, q.table = modeInsert, table
	return q
}

func (q *Q) Select(fields ...string) *Q {
	if len(fields) == 1 {
		split := strings.Split(fields[0], `,`)
		for idx, str := range split {
			split[idx] = strings.TrimSpace(str)
		}
		q.fields = split
	}
	q.fields = fields
	return q
}

func (q *Q) From(table string) *Q {
	q.mode, q.table = modeSelect, table
	return q
}

func (q *Q) Where(query string, value interface{}) *Q {
	q.wheres, q.values = append(q.wheres, query), append(q.values, value)
	return q
}

func (q *Q) Limit(number ...int) *Q {
	if len(number) == 1 {
		q.count = number[0]
	} else if len(number) == 2 {
		q.count, q.count = number[0], number[1]
	}
	return q
}

//Output used by sqlserver
func (q *Q) Output(keys ...string) *Q {
	if len(keys) == 1 {
		split := strings.Split(keys[0], `,`)
		for idx, str := range split {
			split[idx] = strings.TrimSpace(str)
		}
		keys = split
	}
	q.keys = keys
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
	case modeSelect:
		s.WriteString(fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(q.fields, `, `), q.table))
		if len(q.wheres) > 0 {
			s.WriteString(fmt.Sprintf(` WHERE %s`, strings.Join(q.wheres, ` AND `)))
		}
		if q.count > 0 {
			s.WriteString(fmt.Sprintf(` LIMIT %d, %d`, q.start, q.count))
		}
	}
	return s.String()
}

func Query(driverName string) *Q {
	return &Q{driverName: driverName}
}
