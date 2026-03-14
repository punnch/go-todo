package apperrors

import "errors"

var ErrNotFound error = errors.New("task not found")
