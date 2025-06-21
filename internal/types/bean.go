package types

import "reflect"

type Bean interface{}

type LazyBean[T Bean] func() T

type BeanContainer struct {
	Bean          Bean
	Type          reflect.Type
	HoldByContext bool
}

type BeanProvider struct {
	Builder         func(ctx Context) (Bean, error)
	BeanType        reflect.Type
	Options         []func(provider *BeanProvider)
	UseExistingBean reflect.Type
	HoldByContext   bool
}

type ValueBox[T interface{}] struct {
	Value T
}

func EmptyBox[T Bean]() *ValueBox[T] {
	return &ValueBox[T]{}
}
