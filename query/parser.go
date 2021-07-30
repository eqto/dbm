package query

import (
	"regexp"
	"strconv"
	"strings"
)

type Statement interface {
}

type statement struct {
	Statement
	stmt interface{}
}

func Parse(query string) Statement {
	query = strings.TrimSpace(query)
	if strings.HasPrefix(strings.ToUpper(query), `SELECT`) {
		regex := regexp.MustCompile(`(?Uis)^SELECT\s+(?P<fields>.*)\s+FROM\s+(?P<tables>.*)(?:\s+WHERE\s+(?P<where>.*)|)(?:\s+GROUP BY\s+(?P<group>.*)|)(?:\s+ORDER\s+BY\s+(?P<order>.*)|)(?:\s+LIMIT\s+(?P<limit>.+)|)$`)
		matches := regex.FindStringSubmatch(query)
		if len(matches) == 0 {
			return nil
		}
		names := regex.SubexpNames()
		var selectStmt *SelectStmt
		var tableStmt *TableStmt

		for i, name := range names {
			match := matches[i]
			switch name {
			case `fields`:
				selectStmt = Select(match)
			case `tables`:
				tablesRegex := regexp.MustCompile(`(?Uis)\s+(INNER|LEFT|RIGHT)\s+JOIN\s+`)
				tables := tablesRegex.Split(match, -1)
				delimiters := tablesRegex.FindAllString(match, -1)

				delimiterRegex := regexp.MustCompile(`(?Uis)(INNER|LEFT|RIGHT)`)
				joinRegex := regexp.MustCompile(`(?Uis)\s+ON\s+`)
				for i, table := range tables {
					if i == 0 {
						tableStmt = selectStmt.From(table)
					} else {
						joinKind := strings.ToUpper(delimiterRegex.FindString(delimiters[i-1]))
						split := joinRegex.Split(table, 2)
						if len(split) == 2 {
							tableStmt.join(joinKind, strings.TrimSpace(split[0]), strings.TrimSpace(split[1]))
						}
					}
				}
			case `where`:
				whereRegex := regexp.MustCompile(`(?Uis)\s+(AND|OR)\s+`)
				wheres := whereRegex.Split(match, -1)
				delimiters := whereRegex.FindAllString(match, -1)
				var whereStmt *WhereStmt
				for i, where := range wheres {
					if i == 0 {
						whereStmt = tableStmt.Where(strings.TrimSpace(where))
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
			case `group`:
				tableStmt.GroupBy(match)
			case `order`:
				tableStmt.OrderBy(match)
			case `limit`:
				split := strings.Split(match, `,`)
				switch len(split) {
				case 1:
					if count, e := strconv.Atoi(match); e == nil {
						tableStmt.Limit(count)
					}
				case 2:
					if offset, e := strconv.Atoi(strings.TrimSpace(split[0])); e == nil {
						if count, e := strconv.Atoi(strings.TrimSpace(split[1])); e == nil {
							tableStmt.Limit(offset, count)
						}
					}
				}
			}
		}
		return &statement{stmt: selectStmt}
	}
	return nil
}
