package lang

import errors "marketplace-system/lib/customerrors"

var (
	ErrDataNotFound       = errors.NewNotFoundError("Data not found")
	ErrInt64Convert       = errors.NewInternalError("user_id is not of type int64")
	ErrInvalidJsonPayload = errors.NewBadRequestError("Invalid JSON payload")
)
