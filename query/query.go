package query

import (
	"fmt"
	"strings"
)

const (
	ModeSelect = iota
	ModeInsert
	ModeDelete
)

type Builder struct {
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

func (b *Builder) InsertInto(table string, fields ...string) *Builder {
	b.mode, b.table = ModeInsert, table
	return b
}

func (q *Builder) Select(fields ...string) *Builder {
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

func (q *Builder) From(table string) *Builder {
	q.mode, q.table = ModeSelect, table
	return q
}

func (q *Builder) Where(query string, value interface{}) *Builder {
	q.wheres, q.values = append(q.wheres, query), append(q.values, value)
	return q
}

func (q *Builder) Limit(number ...int) *Builder {
	if len(number) == 1 {
		q.count = number[0]
	} else if len(number) == 2 {
		q.count, q.count = number[0], number[1]
	}
	return q
}

//Output used by sqlserver
func (q *Builder) Output(keys ...string) *Builder {
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

func (q *Builder) ValueMap(values map[string]interface{}) *Builder {
	for key, val := range values {
		q.keys = append(q.keys, key)
		q.values = append(q.values, val)
	}
	return q
}

func (q *Builder) Value(key string, value interface{}) *Builder {
	q.keys = append(q.keys, key)
	q.values = append(q.values, value)
	return q
}

func (q *Builder) String() string {
	s := strings.Builder{}
	switch q.mode {
	case ModeInsert:
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
					case ModeInsert:
						outputs = append(outputs, `INSERTED.`+output)
					case ModeDelete:
						outputs = append(outputs, `DELETED.`+output)
					}
				}
				s.WriteString(fmt.Sprintf(` OUTPUT %s`, strings.Join(outputs, `, `)))
			}
		}
		s.WriteString(fmt.Sprintf(` VALUES(%s)`, strings.Join(values, `, `)))
	case ModeSelect:
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

func Build(driverName string) *Builder {
	return &Builder{driverName: driverName}
}
