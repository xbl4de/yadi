package yadi

import (
	"github.com/pkg/errors"
	"log"
	"reflect"
)

var globalCtx Context

var deferredUpdates []func(ctx Context) error

const dummyInt = 42

func SetBeanProvider[T Bean](builder func(ctx Context) (T, error), options ...func(provider *BeanProvider)) int {
	beanType := reflect.TypeFor[T]()
	provider := &BeanProvider{
		builder: func(ctx Context) (Bean, error) {
			return builder(ctx)
		},
		beanType:      beanType,
		holdByContext: true,
	}
	for _, option := range options {
		option(provider)
	}

	provideDefault[T](provider)

	return dummyInt
}

func ProvideAsExistingBean[T, E Bean]() int {
	// TODO: add option
	return dummyInt
}

func WithHoldByUser() func(provider *BeanProvider) {
	return func(provider *BeanProvider) {
		provider.holdByContext = false
	}
}

type BeanProviderContainer struct {
	Provider BeanProvider
	Type     reflect.Type
}

func GetBean[T Bean]() (*ValueBox[T], error) {
	err := ensureContext()
	if err != nil {
		return nil, err
	}

	p := reflect.TypeFor[T]()
	bean, err := globalCtx.Get(p)
	if err != nil {
		return nil, err
	}
	casted, ok := bean.(T)
	if !ok {
		return nil, errors.Errorf("Failed to cast bean to type %s: actual type is %s",
			p.String(), reflect.TypeOf(bean).String())
	}
	return &ValueBox[T]{casted}, nil
}

func RequireBean[T Bean]() T {
	err := ensureContext()
	if err != nil {
		panic(err)
	}

	bean, err := GetBean[T]()
	if err != nil {
		_log.Printf("%+v", err)
		panic(err)
	}
	return bean.Value
}

func ensureContext() error {
	if globalCtx == nil {
		return ErrNilContext
	}
	return nil
}

func GetBeanOrDefault[T Bean](defaultValue T) *ValueBox[T] {
	c, err := GetBean[T]()
	if err != nil {
		return &ValueBox[T]{defaultValue}
	}
	return c
}

func provideDefault[T Bean](provider *BeanProvider) {
	if globalCtx != nil {
		err := globalCtx.Register(provider)
		if err != nil {
			panic(err)
		}
	} else {
		deferredUpdates = append(deferredUpdates, func(ctx Context) error {
			return ctx.Register(provider)
		})
	}
	log.Printf("Provided default bean: %s", reflect.TypeFor[T]().String())
}

func GetValue[T interface{}](path string) (*ValueBox[T], error) {
	err := ensureContext()
	if err != nil {
		return nil, err
	}
	val, err := globalCtx.GetGenericValue(path)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to get value by path: %s", path)
	}
	casted, ok := val.(T)
	if !ok {
		typeName := reflect.TypeFor[T]().String()
		return EmptyBox[T](), errors.Errorf("expected type %s but got %T", typeName, val)
	}
	return &ValueBox[T]{casted}, nil
}

func GetValueOrDefault[T interface{}](path string, defaultValue T) *ValueBox[T] {
	val, err := GetValue[T](path)
	if err != nil {
		return &ValueBox[T]{defaultValue}
	}
	return val
}

func SetValue[T interface{}](path string, value T) int {
	if globalCtx != nil {
		globalCtx.SetGenericValue(path, value)
	} else {
		deferredUpdates = append(deferredUpdates, func(ctx Context) error {
			ctx.SetGenericValue(path, value)
			return nil
		})
	}
	return dummyInt
}

func NewLazyBean[T Bean]() LazyBean[T] {
	return func() T {
		bean, err := GetBean[T]()
		if err != nil {
			panic(err)
		}
		return bean.Value
	}
}

func UseLazyContext() {
	applyContext(NewLazyContext(deferredUpdates))
}

func CloseContext() error {
	err := ensureContext()
	if err != nil {
		return err
	}
	return globalCtx.Close()
}

func applyContext(ctx Context) {
	if globalCtx != nil {
		panic(ErrContextAlreadyExists)
	}
	globalCtx = ctx
}

func resetYadi() {
	deferredUpdates = deferredUpdates[:0]
	if globalCtx != nil {
		err := globalCtx.Close()
		if err != nil {
			panic(err)
		}
	}
	globalCtx = nil
}
