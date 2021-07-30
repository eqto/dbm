package query

type ConditionStmt struct {
	stmt       interface{}
	conditions []string
}

func (c *ConditionStmt) And(query string) *ConditionStmt {
	c.conditions = append(c.conditions, `AND `+query)
	return c
}

func (c *ConditionStmt) Or(query string) *ConditionStmt {
	c.conditions = append(c.conditions, `OR `+query)
	return c
}
