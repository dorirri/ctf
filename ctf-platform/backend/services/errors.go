package services

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadySolved = errors.New("already solved")
var ErrForbidden = errors.New("forbidden")
