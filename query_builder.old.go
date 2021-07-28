package db

// //QueryBuilder ...
// type QueryBuilder struct {
// 	selectOptions string
// 	fields        []field
// 	//field, alias
// 	aliasMap map[string]string

// 	driverKind  string
// 	fromParams  string
// 	whereParams []string
// 	groupParams []string
// 	orderParams []string
// 	limitStart  int
// 	limitLength int
// }

// //Clone ...
// func (q *QueryBuilder) Clone() *QueryBuilder {
// 	clone := &QueryBuilder{
// 		fromParams:  q.fromParams,
// 		limitStart:  q.limitStart,
// 		limitLength: q.limitLength,
// 	}
// 	if q.fields != nil {
// 		clone.fields = make([]field, len(q.fields))
// 		copy(clone.fields, q.fields)
// 	}
// 	if q.aliasMap != nil {
// 		clone.aliasMap = make(map[string]string)
// 		for key, val := range q.aliasMap {
// 			clone.aliasMap[key] = val
// 		}
// 	}

// 	if q.whereParams != nil {
// 		clone.whereParams = make([]string, len(q.whereParams))
// 		copy(clone.whereParams, q.whereParams)
// 	}
// 	if q.groupParams != nil {
// 		clone.groupParams = make([]string, len(q.groupParams))
// 		copy(clone.groupParams, q.groupParams)
// 	}
// 	if q.orderParams != nil {
// 		clone.orderParams = make([]string, len(q.orderParams))
// 		copy(clone.orderParams, q.orderParams)
// 	}
// 	return clone
// }

// func (q *QueryBuilder) parseWhere(where string) {
// 	regex := regexp.MustCompile(`(?is)\s+AND\s+`)
// 	wheres := regex.Split(where, -1)

// 	for _, val := range wheres {
// 		if val != `` {
// 			q.whereParams = append(q.whereParams, val)
// 		}
// 	}
// }

// func (q *QueryBuilder) parseGroup(group string) {
// 	regex := regexp.MustCompile(`(?is)([a-z0-9._]+)(?:\s+([^\s]+)|)\s*(?:,|$)`)
// 	groups := regex.FindAllStringSubmatch(group, -1)

// 	for _, val := range groups {
// 		if val[1] != `` {
// 			order := strings.ToUpper(val[2])
// 			if order == `` {
// 				order = `ASC`
// 			}
// 			q.groupParams = append(q.groupParams, val[1]+` `+order)
// 		}
// 	}
// }

// func (q *QueryBuilder) parseOrder(order string) {
// 	sorts := strings.Split(order, `,`)
// 	for _, val := range sorts {
// 		if val != `` {
// 			q.orderParams = append(q.orderParams, val)
// 		}
// 	}
// }

// func (q *QueryBuilder) parseFields(rawColumns string) {
// 	var buffer bytes.Buffer
// 	var fields []string
// 	sql := rawColumns
// 	if q.aliasMap == nil {
// 		q.aliasMap = make(map[string]string)
// 	}
// 	for {
// 		idx := strings.Index(sql, `,`)
// 		if idx < 0 {
// 			if len(sql) > 0 {
// 				buff := strings.Trim(sql, " \r\n\t")
// 				fields = append(fields, buffer.String()+buff)
// 			}
// 			break
// 		}
// 		buffer.WriteString(sql[0:idx])
// 		buff := strings.Trim(buffer.String(), " \r\n\t")
// 		if strings.Count(buff, `(`) == strings.Count(buff, `)`) {
// 			if len(buff) > 0 {
// 				fields = append(fields, buff)
// 				buffer.Reset()
// 			}
// 		} else {
// 			buffer.WriteString(`, `)
// 		}
// 		sql = sql[idx+1:]
// 	}

// 	regex := regexp.MustCompile(`(?Uis)^(.*)(?:\s+AS\s+(.*)|)$`)

// 	for _, val := range fields {
// 		trimmed := strings.Trim(val, "\r\n\t")
// 		matches := regex.FindStringSubmatch(trimmed)
// 		field := field{name: matches[1], alias: matches[2]}
// 		if matches[2] != `` {
// 			q.aliasMap[matches[2]] = matches[1]
// 		} else {
// 			matches = strings.Split(trimmed, `.`)
// 			if len(matches) > 1 && matches[1] != `` {
// 				q.aliasMap[matches[1]] = trimmed
// 			}
// 		}
// 		q.fields = append(q.fields, field)
// 	}
// }

// //GetField ...
// func (q *QueryBuilder) GetField(name string) string {
// 	if field, ok := q.aliasMap[name]; ok {
// 		return field
// 	}
// 	for _, val := range q.fields {
// 		if val.name == name || strings.HasSuffix(val.name, `.`+name) {
// 			return val.name
// 		}
// 	}
// 	return ``
// }

// //Where ...
// func (q *QueryBuilder) Where(name string) {
// 	q.WhereOp(name, ` = `)
// }

// //Order ...
// func (q *QueryBuilder) Order(field string, order string) {
// 	q.orderParams = append(q.orderParams, field+` `+order)
// }

// //Limit ...
// func (q *QueryBuilder) Limit(start int, length int) {
// 	q.limitLength = length
// 	q.limitStart = start
// }

// //LimitStart ...
// func (q *QueryBuilder) LimitStart() int {
// 	return q.limitStart
// }

// //LimitLength ...
// func (q *QueryBuilder) LimitLength() int {
// 	return q.limitLength
// }

// //WhereOp ...
// func (q *QueryBuilder) WhereOp(name string, operator string) {
// 	if field, ok := q.aliasMap[name]; ok {
// 		name = field
// 	}
// 	if operator == `` {
// 		operator = `=`
// 	}
// 	operator = strings.ToUpper(strings.TrimSpace(operator))
// 	if operator == `LIKE` {
// 		operator = ` LIKE `
// 	}
// 	where := name + operator + `?`

// 	if operator == `FULLTEXT` {
// 		where = `MATCH(` + name + `) AGAINST(? IN BOOLEAN MODE)`
// 	}

// 	q.whereParams = append(q.whereParams, where)
// }

// //ToConditionSQL ...
// func (q *QueryBuilder) ToConditionSQL() string {
// 	var buffer bytes.Buffer
// 	if len(q.whereParams) > 0 {
// 		buffer.WriteString(` WHERE ` + strings.Join(q.whereParams, ` AND `))
// 	}
// 	if len(q.groupParams) > 0 {
// 		buffer.WriteString(` GROUP BY ` + strings.Join(q.groupParams, `, `))
// 	}
// 	if len(q.orderParams) > 0 {
// 		buffer.WriteString(` ORDER BY ` + strings.Join(q.orderParams, `, `))
// 	}
// 	if q.driverKind == `mysql` && q.limitLength > 0 {
// 		buffer.WriteString(` LIMIT ` + strconv.Itoa(q.limitStart) + `, ` + strconv.Itoa(q.limitLength))
// 	}

// 	return buffer.String()
// }

// //ToFromSQL ...
// func (q *QueryBuilder) ToFromSQL() string {
// 	return ` FROM ` + q.fromParams
// }

// //ToSQL ...
// func (q *QueryBuilder) ToSQL() string {
// 	var sqlFields []string
// 	for _, val := range q.fields {
// 		sqlFields = append(sqlFields, val.String())
// 	}

// 	strFields := strings.Join(sqlFields, `, `)
// 	var buffer bytes.Buffer
// 	buffer.WriteString(`SELECT ` + q.selectOptions + strFields + q.ToFromSQL() + q.ToConditionSQL())

// 	return buffer.String()
// }

// //ParseQuery ...
// func ParseQuery(query string) *QueryBuilder {
// 	query = strings.TrimSpace(query)
// 	qb := QueryBuilder{driverKind: `mysql`}
// 	if strings.HasPrefix(strings.ToUpper(query), `SELECT`) {
// 		regex := regexp.MustCompile(`(?Uis)^SELECT\s+(SQL_CALC_FOUND_ROWS\s+|)(.*)\s+FROM\s+(.*)(?:\s+WHERE\s+(.*)|)(?:\s+GROUP BY\s+(.*)|)(?:\s+ORDER\s+BY\s+(.*)|)(?:\s+LIMIT\s+(?:(?:([0-9]+)\s*,\s*|)([0-9]+))|)$`)
// 		matches := regex.FindStringSubmatch(query)

// 		if len(matches) < 4 {
// 			return nil
// 		}
// 		qb.selectOptions = matches[1]
// 		qb.parseFields(matches[2])
// 		qb.fromParams = matches[3]
// 		qb.parseWhere(matches[4])
// 		qb.parseGroup(matches[5])
// 		qb.parseOrder(matches[6])
// 		if matches[6] != `` {
// 			qb.limitStart, _ = strconv.Atoi(matches[6])
// 		}
// 		if matches[7] != `` {
// 			qb.limitLength, _ = strconv.Atoi(matches[7])
// 		}
// 	}

// 	return &qb
// }
