package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateUserName  = errors.New("models: duplicate username")
	ErrFetchingData       = errors.New("models: error fetching data from gutendex")
)
