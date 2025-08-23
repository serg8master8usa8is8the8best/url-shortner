package storage

import "errors"

var (
	ErrURLNotFund = errors.New("url not found")
	ErrURLExists  = errors.New("url exists")
)
