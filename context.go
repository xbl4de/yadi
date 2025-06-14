package yadi

import (
	"reflect"
)

type Context interface {
	Init()
	Close() error
	Register(ctx *BeanProvider) error
	Get(p reflect.Type) (Bean, error)
	GetGenericValue(path string) (interface{}, error)
	SetGenericValue(path string, value interface{})
}
