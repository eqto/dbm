package query

type Where struct {
	conditions []string
}

func (w *Where) And(query string) *Where {
	w.conditions = append(w.conditions, query)
	return w
}
