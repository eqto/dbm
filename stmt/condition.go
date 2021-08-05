package stmt

type Condition struct {
	stmt       interface{}
	conditions []string
}

func (c *Condition) And(query string) *Condition {
	c.conditions = append(c.conditions, `AND `+query)
	return c
}

func (c *Condition) Or(query string) *Condition {
	c.conditions = append(c.conditions, `OR `+query)
	return c
}
