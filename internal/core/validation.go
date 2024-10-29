package core

import "errors"

var ErrValidation = errors.New("validation error")
var ErrNotFound = errors.New("resource not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrMalformedRequest = errors.New("malformed request")