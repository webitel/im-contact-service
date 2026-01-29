package queries

import (
	"unicode/utf8"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type ContactUpdateQuery struct {
	builder sq.UpdateBuilder
}

func NewContactUpdateQuery() *ContactUpdateQuery {
	return &ContactUpdateQuery{
		builder: sq.StatementBuilder.Update(ContactTable).PlaceholderFormat(sq.Dollar),
	}
}

func (q *ContactUpdateQuery) WithDomainIDFilter(domainID int) *ContactUpdateQuery {
	q.builder = q.builder.Where(sq.Eq{"domain_id": domainID})

	return q
}

func (q *ContactUpdateQuery) WithIDFilter(id uuid.UUID) *ContactUpdateQuery {
	q.builder = q.builder.Where(sq.Eq{"id": id})
	
	return q
}

func (q *ContactUpdateQuery) WithName(name string) *ContactUpdateQuery {
	if name != "" && utf8.ValidString(name) {
		q.builder = q.builder.Set("name", name)
	} 
	
	return q
}

func (q *ContactUpdateQuery) WithUsername(username string) *ContactUpdateQuery {
	if username	!= "" && utf8.ValidString(username) {
		q.builder = q.builder.Set("username", username)
	}

	return q
}

func (q *ContactUpdateQuery) WithMetadata(md map[string]string) *ContactUpdateQuery {
	q.builder = q.builder.Set("metadata", md)

	return q
}

func (q *ContactUpdateQuery) WithSubject(sub string) *ContactUpdateQuery {
	if sub != "" && utf8.ValidString(sub) {
		q.builder = q.builder.Set("subject_id", sub)
	}
	
	return q
}

func (q *ContactUpdateQuery) ToSql() (string, []any, error) {
	return q.builder.Suffix("RETURNING *").ToSql()
}