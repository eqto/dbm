package stmt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Parse(s string) interface{} {
	s = strings.TrimSpace(s)
	if len(s) < 6 {
		return nil
	}
	command := s[:6]

	switch strings.ToUpper(command) {
	case `SELECT`:
		RegexStr := `SELECT\s+(?:(?P<distinct>DISTINCT|)\s+|)(?P<fields>.+)\s+FROM\s+(?P<tables>.+)` + `(?:\s+WHERE\s+(?P<where>.+)|)` + `(?:\s+GROUP\s+BY\s+(?P<groupby>.+)|)` + `(?:\s+ORDER\s+BY\s+(?P<orderby>.+)|)` + `(?:\s+LIMIT\s+(?:(?P<offset>[0-9]+)\s*[,]\s*(?P<count>[0-9]+)|(?P<count>[0-9]+))|\s+OFFSET\s+(?P<offset>[0-9]+)\s+ROWS\s+FETCH\s+NEXT\s+(?P<count>[0-9]+)\s+ROWS\s+ONLY|)`

		regex := regexp.MustCompile(fmt.Sprintf(`(?Uis)^\s*(?:%s)\s*$`, RegexStr))

		matches := regex.FindStringSubmatch(s)

		if len(matches) == 0 {
			return nil
		}
		names := regex.SubexpNames()
		var fields *SelectFields
		var stmt *Select

		for i, name := range names {
			match := matches[i]
			if match == `` {
				continue
			}
			switch name {
			case `fields`:
				fields = Build().Select(match)
			case `tables`:
				tablesRegex := regexp.MustCompile(`(?Uis)\s+(INNER|LEFT|RIGHT)\s+JOIN\s+`)
				tables := tablesRegex.Split(match, -1)
				delimiters := tablesRegex.FindAllString(match, -1)

				delimiterRegex := regexp.MustCompile(`(?Uis)(INNER|LEFT|RIGHT)`)
				joinRegex := regexp.MustCompile(`(?Uis)\s+ON\s+`)
				for i, table := range tables {
					if i == 0 {
						stmt = fields.From(table)
					} else {
						joinKind := strings.ToUpper(delimiterRegex.FindString(delimiters[i-1]))
						split := joinRegex.Split(table, 2)
						if len(split) == 2 {
							stmt.join(joinFromString(joinKind), strings.TrimSpace(split[0]), strings.TrimSpace(split[1]))
						}
					}
				}
			case `where`:
				whereRegex := regexp.MustCompile(`(?Uis)\s+(AND|OR)\s+`)
				wheres := whereRegex.Split(match, -1)
				delimiters := whereRegex.FindAllString(match, -1)
				var whereStmt *SelectWhere
				for i, where := range wheres {
					if i == 0 {
						whereStmt = stmt.Where(strings.TrimSpace(where))
					} else {
						delimiter := strings.ToUpper(strings.TrimSpace(delimiters[i-1]))
						switch delimiter {
						case `AND`:
							whereStmt.And(where)
						case `OR`:
							whereStmt.Or(where)
						}
					}
				}
			case `groupby`:
				stmt.GroupBy(match)
			case `orderby`:
				stmt.OrderBy(match)
			case `offset`:
				if i, e := strconv.Atoi(match); e == nil {
					stmt.Offset(i)
				}
			case `count`:
				if i, e := strconv.Atoi(match); e == nil {
					stmt.Count(i)
				}
			}
		}
		return stmt

	case `INSERT`:
	case `UPDATE`:
	case `DELETE`:
	}
	return nil
}

func joinFromString(str string) int {
	switch strings.ToUpper(str) {
	case `LEFT`:
		return LeftJoin
	case `RIGHT`:
		return RightJoin
	case `INNER`:
		fallthrough
	default:
		return InnerJoin
	}
}
