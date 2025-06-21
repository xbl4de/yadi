package di

import (
	"github.com/pkg/errors"
	log2 "github.com/xbl4de/yadi/internal/log"
	"github.com/xbl4de/yadi/internal/utils"
	types2 "github.com/xbl4de/yadi/types"
	"log"
	"reflect"
)

var globalCtx types2.Context

var deferredUpdates []func(ctx types2.Context) error

const dummyInt = 42

func GlobalCtx() types2.Context {
	return globalCtx
}

func SetBeanProvider[T types2.Bean](builder func(ctx types2.Context) (T, error), options ...func(provider *types2.BeanProvider)) int {
	beanType := reflect.TypeFor[T]()
	provider := &types2.BeanProvider{
		Builder: func(ctx types2.Context) (types2.Bean, error) {
			return builder(ctx)
		},
		BeanType:      beanType,
		HoldByContext: true,
	}
	for _, option := range options {
		option(provider)
	}

	provideDefault[T](provider)

	return dummyInt
}

func ProvideAsExistingBean[T, E types2.Bean]() int {
	// TODO: add option
	return dummyInt
}

func WithHoldByUser() func(provider *types2.BeanProvider) {
	return func(provider *types2.BeanProvider) {
		provider.HoldByContext = false
	}
}

type BeanProviderContainer struct {
	Provider types2.BeanProvider
	Type     reflect.Type
}

func GetBean[T types2.Bean]() (*types2.ValueBox[T], error) {
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
	return &types2.ValueBox[T]{Value: casted}, nil
}

func RequireBean[T types2.Bean]() T {
	err := ensureContext()
	if err != nil {
		panic(err)
	}

	bean, err := GetBean[T]()
	if err != nil {
		log2.Log("%+v", err)
		panic(err)
	}
	return bean.Value
}

func ensureContext() error {
	if globalCtx == nil {
		return types2.ErrNilContext
	}
	return nil
}

func GetBeanOrDefault[T types2.Bean](defaultValue T) *types2.ValueBox[T] {
	c, err := GetBean[T]()
	if err != nil {
		return &types2.ValueBox[T]{defaultValue}
	}
	return c
}

func provideDefault[T types2.Bean](provider *types2.BeanProvider) {
	if globalCtx != nil {
		err := globalCtx.Register(provider)
		if err != nil {
			panic(err)
		}
	} else {
		deferredUpdates = append(deferredUpdates, func(ctx types2.Context) error {
			return ctx.Register(provider)
		})
	}
	log.Printf("Provided default bean: %s", reflect.TypeFor[T]().String())
}

func GetValue[T interface{}](path string) (*types2.ValueBox[T], error) {
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
		return types2.EmptyBox[T](), errors.Errorf("expected type %s but got %T", typeName, val)
	}
	return &types2.ValueBox[T]{Value: casted}, nil
}

func GetValueOrDefault[T interface{}](path string, defaultValue T) *types2.ValueBox[T] {
	val, err := GetValue[T](path)
	if err != nil {
		return &types2.ValueBox[T]{Value: defaultValue}
	}
	return val
}

func SetValue[T interface{}](path string, value T) int {
	if globalCtx != nil {
		globalCtx.SetGenericValue(path, value)
	} else {
		deferredUpdates = append(deferredUpdates, func(ctx types2.Context) error {
			ctx.SetGenericValue(path, value)
			return nil
		})
	}
	return dummyInt
}

func NewLazyBean[T types2.Bean]() types2.LazyBean[T] {
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

func ClearDeferredUpdates() {
	deferredUpdates = deferredUpdates[:0]
}

func CloseContext() error {
	err := ensureContext()
	if err != nil {
		return err
	}
	err = globalCtx.Close()
	if err != nil {
		return err
	}
	globalCtx = nil
	return nil
}

func CloseContextSoft() error {
	err := ensureContext()
	if err != nil {
		return nil
	}
	err = globalCtx.Close()
	if err != nil {
		return err
	}
	globalCtx = nil
	return nil
}

func applyContext(ctx types2.Context) {
	if globalCtx != nil {
		panic(types2.ErrContextAlreadyExists)
	}
	globalCtx = ctx
}

func GetBeanFromContext(beanType reflect.Type) (types2.Bean, error) {
	val, err := GlobalCtx().Get(beanType)
	if err != nil {
		return nil, err
	}
	return val, err
}

func GetGenericValueOrDefault(path string, defaultValue interface{}) (interface{}, error) {
	value, err := GlobalCtx().GetGenericValue(path)
	if err != nil {
		if errors.Is(err, types2.ErrNoValueFound) && defaultValue != nil {
			return defaultValue, nil
		}
		return nil, err
	}
	return value, nil
}

func GetGenericValue(path string) (interface{}, error) {
	return GetGenericValueOrDefault(path, nil)
}

func GetBeanOrDefaultFromContext(beanType reflect.Type, defaultValue types2.Bean) (types2.Bean, error) {
	err := utils.ValidateTypeIsBean(beanType)
	if err != nil {
		return nil, err
	}
	bean, err := GetBeanFromContext(beanType)
	if err != nil {
		if types2.ErrNoInjectableProvided(err) && defaultValue != nil {
			return defaultValue, nil
		}
		return nil, err
	}
	return bean, nil
}
