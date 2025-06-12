package yadi

import (
	"errors"
	"reflect"
)

var ErrNoBeanProvider = errors.New("no bean provider found")
var ErrNoValueFound = errors.New("no value found")

type Context interface {
	Init()
	Close() error
	Register(ctx *BeanProvider) error
	Get(p reflect.Type) (Bean, error)
	GetGenericValue(path string) (interface{}, error)
	SetGenericValue(path string, value interface{})
}
