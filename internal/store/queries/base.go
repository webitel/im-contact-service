package queries

type Query interface {
	ToSql() (string, []any, error)
}