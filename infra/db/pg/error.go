package pg

import (
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"

	"github.com/webitel/webitel-go-kit/pkg/errors"
)

// ErrorIntegrityViolation checks if the error is a PostgreSQL integrity constraint violation.
// It maps specific Postgres error codes to their corresponding gRPC status codes
// using the webitel-go-kit error package.
//
// Supported integrity violations:
//   - 23000: integrity_constraint_violation
//   - 23001: restrict_violation
//   - 23502: not_null_violation
//   - 23503: foreign_key_violation
//   - 23505: unique_violation (returns codes.AlreadyExists)
//   - 23514: check_violation
//   - 23P01: exclusion_violation
//
// If the error is a matched integrity violation, it returns the formatted error and true.
// Otherwise, it returns nil and false.
func ErrorIntegrityViolation(err error) (bool, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return true, errors.New(
				"record already exists",
				errors.WithCause(err),
				errors.WithCode(codes.AlreadyExists),
				errors.WithValue("postgresql_code", pgErr.Code),
				errors.WithValue("column", pgErr.ColumnName),
				errors.WithValue("constraint", pgErr.ConstraintName),
				errors.WithValue("detail", pgErr.Detail),
			)
		case "23503", "23502", "23514", "23P01":
			return true, errors.New(
				"integrity constraint violation",
				errors.WithCause(err),
				errors.WithCode(codes.FailedPrecondition),
				errors.WithValue("postgresql_code", pgErr.Code),
				errors.WithValue("constraint", pgErr.ConstraintName),
				errors.WithValue("detail", pgErr.Detail),
			)
		default:
			return false, nil
		}
	}

	return false, nil
}

func ExtractPgErrorMap(err error) map[string]any {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	fields := map[string]any{
		"code":     pgErr.Code,
		"severity": pgErr.Severity,
		"message":  pgErr.Message,
	}

	if pgErr.Detail != "" {
		fields["detail"] = pgErr.Detail
	}

	if pgErr.Hint != "" {
		fields["hint"] = pgErr.Hint
	}

	if pgErr.TableName != "" {
		fields["table"] = pgErr.TableName
	}

	if pgErr.ColumnName != "" {
		fields["column"] = pgErr.ColumnName
	}

	if pgErr.ConstraintName != "" {
		fields["constraint"] = pgErr.ConstraintName
	}

	return fields
}
