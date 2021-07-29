package db

import (
	"strings"
)

const (
	ModeSelect = iota
	ModeInsert
	ModeDelete
)

type QueryBuilder struct {
	drv Driver

	mode         int
	table        Table
	fields       []Field
	keys         []string
	outputs      []string //sqlserver
	values       []interface{}
	wheres       []string
	start, count int
}

func (b *QueryBuilder) InsertInto(table string, fields ...string) *QueryBuilder {
	b.mode = ModeInsert
	b.From(table)
	return b
}

func (q *QueryBuilder) Select(fields ...string) *QueryBuilder {
	if len(fields) == 1 {
		split := strings.Split(fields[0], `,`)
		fields = split
	}
	q.fields = parseFields(fields...)
	return q
}

func parseFields(fields ...string) []Field {
	parsedFields := []Field{}
	for _, field := range fields {
		split := strings.Split(field, ` AS `)
		f := Field{}
		if len(split) == 2 {
			f.alias = strings.TrimSpace(split[1])
			field = split[0]
		}
		split = strings.Split(field, `.`)
		if len(split) == 2 {
			f.table = strings.TrimSpace(split[0])
			f.name = strings.TrimSpace(split[1])
		} else {
			f.name = strings.TrimSpace(field)
		}
		parsedFields = append(parsedFields, f)
	}
	return parsedFields
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
	q.mode = ModeSelect
	split := strings.SplitN(table, ` `, 2)
	if len(split) == 2 {
		q.table = Table{strings.TrimSpace(split[0]), strings.TrimSpace(split[1])}
	} else {
		q.table = Table{strings.TrimSpace(table), ``}
	}
	return q
}

func (q *QueryBuilder) Where(query string, value interface{}) *QueryBuilder {
	q.wheres, q.values = append(q.wheres, query), append(q.values, value)
	return q
}

func (q *QueryBuilder) Limit(number ...int) *QueryBuilder {
	if len(number) == 1 {
		q.count = number[0]
	} else if len(number) == 2 {
		q.count, q.count = number[0], number[1]
	}
	return q
}

//Output used by sqlserver
func (q *QueryBuilder) Output(keys ...string) *QueryBuilder {
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

func (q *QueryBuilder) ValueMap(values map[string]interface{}) *QueryBuilder {
	for key, val := range values {
		q.keys = append(q.keys, key)
		q.values = append(q.values, val)
	}
	return q
}

func (q *QueryBuilder) Value(key string, value interface{}) *QueryBuilder {
	q.keys = append(q.keys, key)
	q.values = append(q.values, value)
	return q
}

func New(drv Driver) *QueryBuilder {
	return &QueryBuilder{drv: drv}
}

// func (q *Builder) String() string {
// 	s := strings.Builder{}
// 	switch q.mode {
// 	case ModeInsert:
// 		values := []string{}
// 		s.WriteString(`INSERT INTO ` + q.table)
// 		if len(q.keys) > 0 {
// 			s.WriteString(`(` + strings.Join(q.keys, `, `) + `)`)
// 			if q.drv.Name() == `sqlserver` {
// 				for idx := range q.values {
// 					values = append(values, fmt.Sprintf(`@p%d`, idx+1))
// 				}
// 			} else {
// 				values = append(values, `?`)
// 			}
// 		}
// 		if len(q.outputs) > 0 {
// 			for _, output := range q.outputs {
// 				ucase := strings.ToUpper(output)
// 				outputs := []string{}
// 				if strings.HasPrefix(ucase, `INSERTED.`) {
// 					outputs = append(outputs, `INSERTED`+output[8:])
// 				} else if strings.HasPrefix(ucase, `DELETED.`) {
// 					outputs = append(outputs, `DELETED`+output[7:])
// 				} else {
// 					switch q.mode {
// 					case ModeInsert:
// 						outputs = append(outputs, `INSERTED.`+output)
// 					case ModeDelete:
// 						outputs = append(outputs, `DELETED.`+output)
// 					}
// 				}
// 				s.WriteString(fmt.Sprintf(` OUTPUT %s`, strings.Join(outputs, `, `)))
// 			}
// 		}
// 		s.WriteString(fmt.Sprintf(` VALUES(%s)`, strings.Join(values, `, `)))
// 	case ModeSelect:
// 		s.WriteString(fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(q.fields, `, `), q.table))
// 		if len(q.wheres) > 0 {
// 			s.WriteString(fmt.Sprintf(` WHERE %s`, strings.Join(q.wheres, ` AND `)))
// 		}
// 		if q.count > 0 {
// 			s.WriteString(fmt.Sprintf(` LIMIT %d, %d`, q.start, q.count))
// 		}
// 	}
// 	return s.String()
// }
