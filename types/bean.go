package types

import "reflect"

type Bean interface{}

type LazyBean[T Bean] func() T

type BeanContainer struct {
	Bean          Bean
	Name          string
	Type          reflect.Type
	HoldByContext bool
}

func NewBeanContainer(
	bean Bean,
	name string,
	typ reflect.Type,
	holdByContext bool,
) *BeanContainer {
	return &BeanContainer{
		Bean:          bean,
		Name:          name,
		Type:          typ,
		HoldByContext: holdByContext,
	}
}

func NewBeanContainerHoldByContext(
	bean Bean,
	name string,
	typ reflect.Type,
) *BeanContainer {
	return &BeanContainer{
		Bean:          bean,
		Name:          name,
		Type:          typ,
		HoldByContext: true,
	}
}

func NewBeanContainerHoldByUser(
	bean Bean,
	name string,
	typ reflect.Type,
) *BeanContainer {
	return &BeanContainer{
		Bean:          bean,
		Name:          name,
		Type:          typ,
		HoldByContext: false,
	}
}

type BeanProvider struct {
	Builder         func(ctx Context) (Bean, error)
	BeanType        reflect.Type
	BeanName        string
	Options         []func(provider *BeanProvider)
	UseExistingBean reflect.Type
	HoldByContext   bool
}
