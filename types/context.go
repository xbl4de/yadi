package types

import (
	"reflect"
)

type Context interface {
	Init()
	Close() error
	Register(ctx *BeanProvider) error
	Get(typ reflect.Type) (Bean, error)
	GetNamed(typ reflect.Type, beanName string) (Bean, error)
	GetGenericValue(path string) (interface{}, error)
	SetGenericValue(path string, value interface{})
}
