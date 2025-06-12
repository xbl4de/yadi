package yadi

import (
	"github.com/pkg/errors"
	"log"
	"reflect"
)

var globalCtx Context

var deferredUpdates []func(ctx Context) error

func SetBeanProvider[T Bean](builder func(ctx Context) (T, error), options ...func(provider *BeanProvider)) int {
	beanType := reflect.TypeFor[T]()
	provider := &BeanProvider{
		builder: func(ctx Context) (Bean, error) {
			return builder(ctx)
		},
		beanType: beanType,
	}
	for _, option := range options {
		option(provider)
	}

	provideDefault[T](provider)

	return 42 // doesn't work with any other value...just joke(?)
}

func SetBeanProviderFunc[T Bean](function interface{}) int {
	return SetBeanProvider(func(ctx Context) (T, error) {
		box, err := providerFromFuncE[T](function)
		return box.Value, err
	})
}

func providerFromFuncE[T Bean](function interface{}) (*ValueBox[T], error) {
	val := reflect.ValueOf(function)
	if val.Kind() != reflect.Func {
		return EmptyBox[T](), errors.Errorf("Expected function, but provided %s", val.Kind().String())
	}
	t := reflect.TypeOf(function)
	if t.NumOut() > 2 {
		return EmptyBox[T](), errors.Errorf("Provider function must return at most two values, but got %d", t.NumIn())
	}
	if t.NumOut() == 2 && !t.Out(1).Implements(reflect.TypeFor[error]()) {
		return EmptyBox[T](), errors.Errorf("Provider function must return an error as second value, but got %s", t.Out(1).Kind().String())
	}
	args := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		argT := t.In(i)
		bean, err := globalCtx.Get(argT)
		if err != nil {
			return EmptyBox[T](), err
		}
		args[i] = reflect.ValueOf(bean)
	}
	outs := val.Call(args)
	box := &ValueBox[T]{outs[0].Interface().(T)}
	if len(outs) == 1 {
		return box, nil
	} else {
		return box, outs[1].Interface().(error)
	}
}

func WithExistingBean[T, E Bean]() int {
	_ = SetBeanProvider(func(ctx Context) (Bean, error) {
		return nil, nil
	}, func(provider *BeanProvider) {
		provider.beanType = reflect.TypeFor[E]()
		provider.useExistingBean = reflect.TypeFor[T]()
	})
	return 1337 // doesn't work with any other value...just joke(?)
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

func NewBeanContext[T Bean](provider BeanProvider) *BeanProviderContainer {
	return &BeanProviderContainer{
		Provider: provider,
		Type:     reflect.TypeFor[T](),
	}
}
func Provide(provider *BeanProvider) error {
	err := globalCtx.Register(provider)
	if err != nil {
		return err
	}
	return nil
}

func GetBean[T Bean]() (*ValueBox[T], error) {
	p := reflect.TypeFor[T]()
	bean, err := globalCtx.Get(p)
	if err != nil {
		return nil, err
	}
	casted, ok := bean.(T)
	if !ok {
		return nil, errors.WithMessagef(err, "Failed to obtain bean from context")
	}
	return &ValueBox[T]{casted}, nil
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
	val, err := globalCtx.GetGenericValue(path)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to get value by path: %s", path)
	}
	casted, ok := val.(T)
	if !ok {
		return nil, errors.WithMessagef(err, "expected type %s but got %T", reflect.TypeFor[T]().String(), val)
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
	return 2
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
	globalCtx = NewLazyContext(deferredUpdates)
}
