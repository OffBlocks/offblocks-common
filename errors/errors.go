package errors

import "errors"

var ErrNotFound = errors.New("resource not found")

var ErrUnsupported = errors.ErrUnsupported

var ErrInvalid = errors.New("validation failed")

var ErrInternal = errors.New("internal error")
