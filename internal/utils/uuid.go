package utils

import (
	"github.com/google/uuid"

	"github.com/webitel/webitel-go-kit/pkg/errors"
)

func ParseStringToUUID(in string, out *uuid.UUID) error {
	if out == nil {
		return errors.InvalidArgument("output pointer is nil", errors.WithID("utils.uuid.parse_string_to_uuid"))
	}

	parsed, err := uuid.Parse(in)
	if err != nil {
		return errors.InvalidArgument(
			"parsing input UUID string",
			errors.WithCause(err),
			errors.WithID("utils.uuid.parse_string_to_uuid"),
		)
	}

	*out = parsed

	return nil
}
