package yadi

import (
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/log"
	"github.com/xbl4de/yadi/types"
	"reflect"
)

func SetBeanProvider[T types.Bean](builder func(ctx types.Context) (T, error), options ...func(provider *types.BeanProvider)) int {
	beanType := reflect.TypeFor[T]()
	provider := &types.BeanProvider{
		Builder: func(ctx types.Context) (types.Bean, error) {
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

func WithHoldByUser() func(provider *types.BeanProvider) {
	return func(provider *types.BeanProvider) {
		provider.HoldByContext = false
	}
}

func WithBeanName(name string) func(provider *types.BeanProvider) {
	return func(provider *types.BeanProvider) {
		provider.BeanName = name
	}
}

func ProvideAsExistingBean[T, E types.Bean]() int {
	panic("not implemented")
}

func GetBean[T types.Bean]() (T, error) {
	err := ensureContext()
	var zeroValue T
	if err != nil {
		return zeroValue, err
	}

	p := reflect.TypeFor[T]()
	bean, err := globalCtx.Get(p)
	if err != nil {
		return zeroValue, err
	}
	casted, ok := bean.(T)
	if !ok {
		return zeroValue, errors.Errorf("Failed to cast bean to type %s: actual type is %s",
			p.String(), reflect.TypeOf(bean).String())
	}
	return casted, nil
}

func GetNamedBean[T types.Bean](name string) (T, error) {
	err := ensureContext()
	var zeroValue T
	if err != nil {
		return zeroValue, err
	}

	typ := reflect.TypeFor[T]()
	bean, err := globalCtx.GetNamed(typ, name)
	if err != nil {
		return zeroValue, err
	}
	casted, ok := bean.(T)
	if !ok {
		return zeroValue, errors.Errorf("Failed to cast bean to type %s: actual type is %s",
			typ.String(), reflect.TypeOf(bean).String())
	}
	return casted, nil
}

func RequireBean[T types.Bean]() T {
	err := ensureContext()
	if err != nil {
		panic(err)
	}

	bean, err := GetBean[T]()
	if err != nil {
		log.Log("%+v", err)
		panic(err)
	}
	return bean
}

func RequireNamedBean[T types.Bean](name string) T {
	err := ensureContext()
	if err != nil {
		panic(err)
	}
	bean, err := GetNamedBean[T](name)
	if err != nil {
		log.Log("%+v", err)
		panic(err)
	}
	return bean
}

func GetBeanOrDefault[T types.Bean](defaultValue T) T {
	bean, err := GetBean[T]()
	if err != nil {
		return defaultValue
	}
	return bean
}

func GetNamedBeanOrDefault[T types.Bean](name string, defaultValue T) T {
	bean, err := GetNamedBean[T](name)
	if err != nil {
		return defaultValue
	}
	return bean
}

func GetValue[T interface{}](path string) (T, error) {
	err := ensureContext()
	var zeroValue T
	if err != nil {
		return zeroValue, err
	}
	val, err := globalCtx.GetGenericValue(path)
	if err != nil {
		return zeroValue, errors.WithMessagef(err, "Failed to get value by path: %s", path)
	}
	casted, ok := val.(T)
	if !ok {
		typeName := reflect.TypeFor[T]().String()
		return zeroValue, errors.Errorf("expected type %s but got %T", typeName, val)
	}
	return casted, nil
}

func GetValueOrDefault[T interface{}](path string, defaultValue T) T {
	val, err := GetValue[T](path)
	if err != nil {
		return defaultValue
	}
	return val
}

func SetValue[T interface{}](path string, value T) int {
	if globalCtx != nil {
		globalCtx.SetGenericValue(path, value)
	} else {
		deferredUpdates = append(deferredUpdates, func(ctx types.Context) error {
			ctx.SetGenericValue(path, value)
			return nil
		})
	}
	return dummyInt
}

func NewLazyBean[T types.Bean]() types.LazyBean[T] {
	return func() T {
		return RequireBean[T]()
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
	err = globalCtx.Close()
	if err != nil {
		return err
	}
	globalCtx = nil
	return nil
}

func WithValuePathAt(paramIndex int, path string) FuncProviderOption {
	return func(config *FuncProviderConfig) {
		config.Parameter(paramIndex).ValuePath = path
	}
}

func WithFuncProviderBeanName(name string) FuncProviderOption {
	return func(config *FuncProviderConfig) {
		config.beanName = name
	}
}

func WithDefaultValueAt(paramIndex int, defaultValue interface{}) FuncProviderOption {
	return func(config *FuncProviderConfig) {
		config.Parameter(paramIndex).DefaultValue = defaultValue
	}
}

func SetBeanProviderFunc[T types.Bean](function interface{}, opts ...FuncProviderOption) int {
	beanName := extractBeanName(opts)
	return SetBeanProvider(func(ctx types.Context) (T, error) {
		cfg := NewFuncProviderConfig()
		for _, opt := range opts {
			opt(cfg)
		}
		bean, err := providerFromFuncE[T](function, cfg)
		return bean, err
	}, WithBeanName(beanName))
}

func InjectLazyBean[T types.Bean]() types.LazyBean[T] {
	return func() T {
		return RequireBean[T]()
	}
}

func Inject(valuePtr types.Bean) error {
	return injectToPtr(reflect.ValueOf(valuePtr))
}
