package erro

import (
	"errors"
)

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrAccountExists    = errors.New("account already exists")
	ErrAlreadySuscribed = errors.New("account already suscribed to this deck")
	ErrNotSuscribed     = errors.New("account isn't suscribed to this deck")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrDeckNotFound     = errors.New("deck not found")
	ErrDeckExists       = errors.New("deck already exists")
	ErrCardNotFound     = errors.New("card not found")
	ErrWrongNotFound    = errors.New("wrong answer not found")
	ErrCardExists       = errors.New("card already exists")

	ErrBadField = errors.New("field is empty or invalid")
)
