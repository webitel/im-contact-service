package queries

type Query interface {
	ToSQL() (string, []any, error)
}
