package stmt

import (
	"strings"
)

type Select struct {
	fields        SelectFields
	tables        Tables
	wheres        []WhereParam
	groupBy       string
	orderBy       string
	offset, count int
}

func (s *Select) InnerJoin(table, on string) *Select {
	return s.join(InnerJoin, table, on)
}

func (s *Select) LeftJoin(table, on string) *Select {
	return s.join(LeftJoin, table, on)
}

func (s *Select) RightJoin(table, on string) *Select {
	return s.join(RightJoin, table, on)
}

func (s *Select) Where(condition string) *SelectWhere {
	s.wheres = []WhereParam{{condition, false}}
	return &SelectWhere{Where{stmt: s}}
}

func (s *Select) GroupBy(groupBy string) *GroupBy {
	s.groupBy = groupBy
	return &GroupBy{s}
}

func (s *Select) OrderBy(orderBy string) *OrderBy {
	s.orderBy = orderBy
	return &OrderBy{s}
}

func (s *Select) Limit(num ...int) *Select {
	switch len(num) {
	case 2:
		s.offset = num[0]
		s.count = num[1]
	case 1:
		s.count = num[0]
	}
	return s
}

func (s *Select) join(joinKind int, table, on string) *Select {
	split := strings.SplitN(table, ` `, 2)
	s.tables.add(strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), joinKind, on)
	return s
}
