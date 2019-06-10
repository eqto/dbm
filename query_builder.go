/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2017-12-17 11:23:27
 */

package db

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

//QueryBuilder ...
type QueryBuilder struct {
	fields []field
	//field, alias
	aliasMap map[string]string
	// fieldMap map[string]string

	whereParams []string
	fromParams  string
	orderParams []string
	limitStart  int
	limitLength int
}

type field struct {
	name, alias string
}

//String ...
func (f *field) String() string {
	if f.alias == `` {
		return f.name
	}
	return f.name + ` AS ` + f.alias
}

//Parse ...
func Parse(query string) *QueryBuilder {
	qb := QueryBuilder{}
	regex := regexp.MustCompile(`(?Uis)^\s*SELECT\s+(.*)\s+FROM\s+(.*)(?:\s+WHERE\s+(.*)|)(?:\s+ORDER\s+BY\s+(.*)|)\s*$`)
	matches := regex.FindStringSubmatch(query)

	sorts := strings.Split(matches[4], `,`)
	for _, val := range sorts {
		if val != `` {
			qb.orderParams = append(qb.orderParams, val)
		}
	}

	regex = regexp.MustCompile(`(?is)\s+AND\s+`)
	wheres := regex.Split(matches[3], -1)

	for _, val := range wheres {
		if val != `` {
			qb.whereParams = append(qb.whereParams, val)
		}
	}

	if len(matches) < 3 {
		return nil
	}
	qb.fromParams = matches[2]
	qb.splitColumns(matches[1])

	return &qb
}

func (q *QueryBuilder) splitColumns(rawColumns string) {
	var buffer bytes.Buffer
	var fields []string
	sql := rawColumns
	if q.aliasMap == nil {
		q.aliasMap = make(map[string]string)
	}
	for {
		idx := strings.Index(sql, `,`)
		if idx < 0 {
			if len(sql) > 0 {
				buff := strings.Trim(sql, " \r\n\t")
				fields = append(fields, buffer.String()+buff)
			}
			break
		}
		buffer.WriteString(sql[0:idx])
		buff := strings.Trim(buffer.String(), " \r\n\t")
		if strings.Count(buff, `(`) == strings.Count(buff, `)`) {
			if len(buff) > 0 {
				fields = append(fields, buff)
				buffer.Reset()
			}
		} else {
			buffer.WriteString(`, `)
		}
		sql = sql[idx+1:]
	}

	regex := regexp.MustCompile(`(?Uis)^(.*)(?:\s+AS\s+(.*)|)$`)

	for _, val := range fields {
		trimmed := strings.Trim(val, "\r\n\t")
		matches := regex.FindStringSubmatch(trimmed)
		field := field{name: matches[1], alias: matches[2]}
		if matches[2] != `` {
			q.aliasMap[matches[2]] = matches[1]
		} else {
			matches = strings.Split(trimmed, `.`)
			if len(matches) > 1 && matches[1] != `` {
				q.aliasMap[matches[1]] = trimmed
			}
		}
		q.fields = append(q.fields, field)
	}
}

//GetField ...
func (q *QueryBuilder) GetField(name string) string {
	if field, ok := q.aliasMap[name]; ok {
		return field
	}
	for _, val := range q.fields {
		if val.name == name || strings.HasSuffix(val.name, `.`+name) {
			return val.name
		}
	}
	return ``
}

//Where ...
func (q *QueryBuilder) Where(name string) {
	q.WhereOp(name, ` = `)
}

//Order ...
func (q *QueryBuilder) Order(field string, order string) {
	q.orderParams = append(q.orderParams, field+` `+order)
}

//Limit ...
func (q *QueryBuilder) Limit(start int, length int) {
	q.limitLength = length
	q.limitStart = start
}

//WhereOp ...
func (q *QueryBuilder) WhereOp(name string, operator string) {
	if field, ok := q.aliasMap[name]; ok {
		name = field
	}
	if operator == `` {
		operator = ` = `
	}
	where := name + operator + `?`
	q.whereParams = append(q.whereParams, where)
}

//ToConditionSQL ...
func (q *QueryBuilder) ToConditionSQL() string {
	var buffer bytes.Buffer
	if len(q.whereParams) > 0 {
		buffer.WriteString(` WHERE ` + strings.Join(q.whereParams, ` AND `))
	}
	if len(q.orderParams) > 0 {
		buffer.WriteString(` ORDER BY ` + strings.Join(q.orderParams, `, `))
	}
	if q.limitLength > 0 {
		buffer.WriteString(` LIMIT ` + strconv.Itoa(q.limitStart) + `, ` + strconv.Itoa(q.limitLength))
	}

	return buffer.String()
}

//ToFromSQL ...
func (q *QueryBuilder) ToFromSQL() string {
	return ` FROM ` + q.fromParams
}

//ToSQL ...
func (q *QueryBuilder) ToSQL() string {
	var sqlFields []string
	for _, val := range q.fields {
		sqlFields = append(sqlFields, val.String())
	}

	strFields := strings.Join(sqlFields, `, `)
	var buffer bytes.Buffer
	buffer.WriteString(`SELECT ` + strFields + q.ToFromSQL() + q.ToConditionSQL())

	return buffer.String()
}
