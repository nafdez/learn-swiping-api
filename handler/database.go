package handle

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
)

type Database interface {
	Create() error
	Read(table, fields, filter string) (struct{}, error)
	Update() error
	Delete() error
}
