package postgres_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"

	"github.com/webitel/im-contact-service/internal/store/postgres"
)

func TestExtractSortingOperator(t *testing.T) {
	testCases := []struct {
		name          string
		inputSort     string
		expectedField string
		expectedDir   string
	}{
		{
			name:          "Given empty string, When extracting, Then returns empty field and direction",
			inputSort:     "",
			expectedField: "",
			expectedDir:   "",
		},
		{
			name:          "Given sort with plus prefix, When extracting, Then returns field and DESC direction",
			inputSort:     "+created_at",
			expectedField: "created_at",
			expectedDir:   postgres.DESC,
		},
		{
			name:          "Given sort with minus prefix, When extracting, Then returns field and ASC direction",
			inputSort:     "-created_at",
			expectedField: "created_at",
			expectedDir:   postgres.ASC,
		},
		{
			name:          "Given sort without prefix, When extracting, Then returns field and default ASC direction",
			inputSort:     "name",
			expectedField: "name",
			expectedDir:   postgres.ASC,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field, dir := postgres.ExtractSortingOperator(tc.inputSort)

			if field != tc.expectedField {
				t.Errorf("Expected field '%s', got '%s'", tc.expectedField, field)
			}

			if dir != tc.expectedDir {
				t.Errorf("Expected direction '%s', got '%s'", tc.expectedDir, dir)
			}
		})
	}
}

func TestApplySortOrDefault(t *testing.T) {
	allowedFields := []string{"id", "name", "created_at"}

	t.Run("Given empty provided sort, When applying sort, Then falls back to default sort", func(t *testing.T) {
		sb := sq.Select("*").From("users")
		defaultSort := "+created_at"

		resultSb, err := postgres.ApplySortOrDefault(allowedFields, "", defaultSort, sb)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		sql, _, err := resultSb.ToSql()
		if err != nil {
			t.Fatalf("Failed to build SQL: %v", err)
		}

		expectedSQL := "SELECT * FROM users ORDER BY created_at desc"
		if sql != expectedSQL {
			t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
		}
	})

	t.Run("Given provided sort with only prefix (empty field), When applying sort, Then falls back to default sort", func(t *testing.T) {
		sb := sq.Select("*").From("users")
		provided := "+"
		defaultSort := "-name"

		resultSb, err := postgres.ApplySortOrDefault(allowedFields, provided, defaultSort, sb)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		sql, _, err := resultSb.ToSql()
		if err != nil {
			t.Fatalf("Failed to build SQL: %v", err)
		}

		expectedSQL := "SELECT * FROM users ORDER BY name asc"
		if sql != expectedSQL {
			t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
		}
	})

	t.Run("Given unsupported sort field, When applying sort, Then returns error", func(t *testing.T) {
		sb := sq.Select("*").From("users")
		provided := "hacker_field"
		defaultSort := "name"

		if _, err := postgres.ApplySortOrDefault(allowedFields, provided, defaultSort, sb); err == nil {
			t.Error("Expected error for unsupported sort field, got nil")
		}
	})

	t.Run("Given valid provided sort, When applying sort, Then applies provided sort successfully", func(t *testing.T) {
		sb := sq.Select("*").From("users")
		provided := "-name"
		defaultSort := "+created_at"

		resultSb, err := postgres.ApplySortOrDefault(allowedFields, provided, defaultSort, sb)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		sql, _, err := resultSb.ToSql()
		if err != nil {
			t.Fatalf("Failed to build SQL: %v", err)
		}

		expectedSQL := "SELECT * FROM users ORDER BY name asc"
		if sql != expectedSQL {
			t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
		}
	})

	t.Run("Given SQL injection to sort.", func(t *testing.T) {
		sb := sq.Select("*").From("users")
		provided := "(case when 1=1 then (select 1 from pg_sleep(5)) else 1 end)"
		defaultSort := "-created_at"

		if _, err := postgres.ApplySortOrDefault(allowedFields, provided, defaultSort, sb); err == nil {
			t.Error("Expected error due to SQL injection, get nil")
		}
	})
}
