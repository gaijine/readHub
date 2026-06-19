package service

import (
	"errors"
)

var (
	ErrBookNotOwned                  = errors.New("book does not belong to user")
	ErrNegativePage                  = errors.New("the number must not be negative")
	ErrPageExceedsBookLength         = errors.New("the number should not be greater than the total number of pages")
	ErrActiveSessionNotFound         = errors.New("you don't have an active session")
	ErrActiveSessionIsExist          = errors.New("an active session already exists")
	ErrBookAlreadyCompleted          = errors.New("book already completed")
	ErrInvalidTotalPages             = errors.New("you entered an incorrect total number of pages, it should not be negative or equal to zero.")
	ErrTotalPagesLessThanCurrentPage = errors.New("the total number of pages cannot be less than the current page")
)
