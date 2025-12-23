package store

import (
	"fmt"
	"slices"
)

const defaultSort string = "created_at desc"

func ValidateAndFormatSort(sort string, allowedFields []string) string {
	if len(sort) < 2 {
		return defaultSort
	}

	var (
		direction = "asc"
		column    = sort[1:]
	)

	if sort[0] == '-' {
		direction = "desc"
	} else if sort[0] != '+' {
		return defaultSort
	}

	if !slices.Contains(allowedFields, column) {
		return defaultSort
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func SanitizeFields(inputFields, allowedFields []string) []string {
	return slices.DeleteFunc(inputFields, func(f string) bool {
		return !slices.Contains(allowedFields, f)
	})
}
