package stmt

import (
	"strings"
)

type Select struct {
	fields        SelectFields
	tables        Tables
	wheres        []WhereParam
	group         string
	order         string
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

func (s *Select) Where(conditions ...string) *SelectWhere {
	if len(conditions) > 0 {
		for _, c := range conditions {
			s.wheres = append(s.wheres, WhereParam{c, false})
		}
	}
	return &SelectWhere{s}
}

func (s *Select) GroupBy(groupBy string) *GroupBy {
	s.group = groupBy
	return &GroupBy{s}
}

func (s *Select) OrderBy(orderBy string) *OrderBy {
	s.order = orderBy
	return &OrderBy{s}
}

func (s *Select) Offset(offset int) *Select {
	s.offset = offset
	return s
}

func (s *Select) Count(count int) *Select {
	s.count = count
	return s
}

func (s *Select) Limit(num ...int) *Select {
	switch len(num) {
	case 2:
		s.Offset(num[0]).Count(num[1])
	case 1:
		s.Count(num[0])
	}
	return s
}

func (s *Select) join(joinKind int, table, on string) *Select {
	split := strings.SplitN(table, ` `, 2)
	s.tables.add(strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), joinKind, on)
	return s
}
