package yadi

import "reflect"

type Bean interface{}

type LazyBean[T Bean] func() T

type BeanContainer struct {
	Bean          Bean
	Type          reflect.Type
	HoldByContext bool
}

type BeanProvider struct {
	builder         func(ctx Context) (Bean, error)
	beanType        reflect.Type
	options         []func(provider *BeanProvider)
	useExistingBean reflect.Type
	holdByContext   bool
}

type ValueBox[T interface{}] struct {
	Value T
}

func EmptyBox[T Bean]() *ValueBox[T] {
	return &ValueBox[T]{}
}
