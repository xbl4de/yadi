package yadi

import (
	"github.com/xbl4de/yadi/internal/di"
	"github.com/xbl4de/yadi/internal/types"
)

func SetBeanProvider[T types.Bean](builder func(ctx types.Context) (T, error), options ...func(provider *types.BeanProvider)) int {
	return di.SetBeanProvider(builder, options...)
}

func WithHoldByUser() func(provider *types.BeanProvider) {
	return di.WithHoldByUser()
}

func GetBean[T types.Bean]() (*types.ValueBox[T], error) {
	return di.GetBean[T]()
}

func RequireBean[T types.Bean]() T {
	return di.RequireBean[T]()
}

func GetBeanOrDefault[T types.Bean](defaultValue T) *types.ValueBox[T] {
	return di.GetBeanOrDefault[T](defaultValue)
}

func GetValue[T interface{}](path string) (*types.ValueBox[T], error) {
	return di.GetValue[T](path)
}

func GetValueOrDefault[T interface{}](path string, defaultValue T) *types.ValueBox[T] {
	return di.GetValueOrDefault[T](path, defaultValue)
}

func SetValue[T interface{}](path string, value T) int {
	return di.SetValue[T](path, value)
}

func NewLazyBean[T types.Bean]() types.LazyBean[T] {
	return di.NewLazyBean[T]()
}

func UseLazyContext() {
	di.UseLazyContext()
}

func CloseContext() error {
	return di.CloseContext()
}
