package store

import "errors"

var (
	ErrUpdatedContactNotFound = errors.New("no result on partial update")
)