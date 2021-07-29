package query

type Where struct {
	stmt       interface{}
	conditions []string
	orderFunc  func(string) *OrderBy
}

func (w *Where) And(query string) *Where {
	w.conditions = append(w.conditions, `AND `+query)
	return w
}

func (w *Where) Or(query string) *Where {
	w.conditions = append(w.conditions, `OR `+query)
	return w
}

func (w *Where) OrderBy(order string) *OrderBy {
	return w.orderFunc(order)
}
