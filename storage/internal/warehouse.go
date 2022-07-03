package internal

import (
	"fmt"
)

type Warehouse interface {
	Get(k string) (string, error)
	Put(k, v string) error
	Delete(k string) error
}

type UnknownKeyError struct {
	k string
}

func NewUnknownKeyError(k string) error {
	return UnknownKeyError{k: k}
}

func (e UnknownKeyError) Error() string {
	return fmt.Sprintf("record with the unknown key was requested: %s", e.k)
}
