package query

type Del struct {
	Key string
}

func (q *Del) GetType() string {
	return DelQuery
}
