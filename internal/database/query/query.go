package query

type QueryInterface interface {
	GetType() string
}

var (
	SetQuery = "SET"
	GetQuery = "GET"
	DelQuery = "DEL"
)
