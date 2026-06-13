package service

import (
	"errors"
)

var (
	ErrBookNotOwned          = errors.New("book does not belong to user")
	ErrNegativePage          = errors.New("the number must not be negative")
	ErrPageExceedsBookLength = errors.New("the number should not be greater than the total number of pages")
	ErrStatusCode            = errors.New
)
