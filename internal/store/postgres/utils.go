package postgres

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const (
	MaxLimit     int = 500
	DefaultLimit int = 20
)

const (
	ASC  string = "asc"
	DESC string = "desc"
)

func ApplyPaging(page, size int, sb sq.SelectBuilder) sq.SelectBuilder {
	if size > 0 {
		if size <= 0 || size > MaxLimit {
			size = DefaultLimit
		}

		sb = sb.Limit(uint64(size + 1))
		if page > 1 {
			sb = sb.Offset(uint64((page - 1) * size))
		}
	}

	return sb
}

func Ident(left, right string) string { return left + "." + right }

func ExtractSortingOperator(sort string) (string, string) {
	if len(sort) != 0 {
		desc := strings.HasPrefix(sort, "+")
		asc := strings.HasPrefix(sort, "-")

		if desc || asc {
			sort = sort[1:]
		}

		dir := ASC
		if desc {
			dir = DESC
		}

		return sort, dir
	}

	return "", ""
}
