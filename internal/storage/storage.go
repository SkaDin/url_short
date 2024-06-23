package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLIsExists = errors.New("url exists")
)
