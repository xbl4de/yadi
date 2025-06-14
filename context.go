package yadi

import (
	"errors"
	"reflect"
)

var ErrNoBeanProvider = errors.New("no bean provider found")
var ErrNoValueFound = errors.New("no value found")

func errNoContextProvided(err error) bool {
	return errors.Is(err, ErrNoBeanProvider) || errors.Is(err, ErrNoValueFound)
}

type Context interface {
	Init()
	Close() error
	Register(ctx *BeanProvider) error
	Get(p reflect.Type) (Bean, error)
	GetGenericValue(path string) (interface{}, error)
	SetGenericValue(path string, value interface{})
}
