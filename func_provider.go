package yadi

import (
	"github.com/pkg/errors"
	"reflect"
)

type ParameterConfig struct {
	ValuePath    string
	DefaultValue interface{}
}

type FuncProviderConfig struct {
	parameters map[int]*ParameterConfig
}

func (c *FuncProviderConfig) Parameter(index int) *ParameterConfig {
	if index < 0 {
		panic(errors.New("index must be >= 0"))
	}
	var p *ParameterConfig
	var ok bool
	if p, ok = c.parameters[index]; !ok {
		p = &ParameterConfig{}
		c.parameters[index] = p
	}
	return p
}

func NewFuncProviderConfig() *FuncProviderConfig {
	return &FuncProviderConfig{
		parameters: make(map[int]*ParameterConfig),
	}
}

type FuncProviderOption func(*FuncProviderConfig)

func WithValuePathAt(paramIndex int, path string) FuncProviderOption {
	return func(config *FuncProviderConfig) {
		config.Parameter(paramIndex).ValuePath = path
	}
}

func WithDefaultValueAt(paramIndex int, defaultValue interface{}) FuncProviderOption {
	return func(config *FuncProviderConfig) {
		config.Parameter(paramIndex).DefaultValue = defaultValue
	}
}

func SetBeanProviderFunc[T Bean](function interface{}, opts ...FuncProviderOption) int {
	return SetBeanProvider(func(ctx Context) (T, error) {
		cfg := NewFuncProviderConfig()
		for _, opt := range opts {
			opt(cfg)
		}
		box, err := providerFromFuncE[T](function, cfg)
		return box.Value, err
	})
}

func providerFromFuncE[T Bean](function interface{}, cfg *FuncProviderConfig) (*ValueBox[T], error) {
	funcValue := reflect.ValueOf(function)
	funcType := reflect.TypeOf(function)

	err := validateProviderFunc[T](funcValue, funcType)
	if err != nil {
		return EmptyBox[T](), err
	}

	args, err := buildArgs(funcType, cfg)
	if err != nil {
		return EmptyBox[T](), err
	}

	outs := funcValue.Call(args)
	box := buildValueBox[T](outs)

	if len(outs) == 1 {
		return box, nil
	} else {
		return box, castToErr(outs[1])
	}
}

func castToErr(errVal reflect.Value) error {
	if errVal.IsNil() {
		return nil
	}
	return errVal.Elem().Interface().(error)
}

func buildValueBox[T interface{}](outs []reflect.Value) *ValueBox[T] {
	var val T
	if outs[0].Kind() == reflect.Interface {
		val = outs[0].Elem().Interface().(T)
	} else {
		val = outs[0].Interface().(T)
	}
	box := &ValueBox[T]{val}
	return box
}

func buildArgs(funcType reflect.Type, cfg *FuncProviderConfig) ([]reflect.Value, error) {
	args := make([]reflect.Value, funcType.NumIn())
	for i := 0; i < funcType.NumIn(); i++ {
		arg, err := findArgValue(funcType.In(i), cfg.Parameter(i))
		if err != nil {
			return nil, errors.WithMessagef(err, "Failed to find arg at index %d", i)
		}
		args[i] = reflect.ValueOf(arg)
	}
	return args, nil
}

func validateProviderFunc[T interface{}](funcValue reflect.Value, funcType reflect.Type) error {
	if funcValue.Kind() != reflect.Func {
		return errors.Errorf("Expected function, but provided %s", funcValue.Kind().String())
	}

	if funcType.NumOut() == 0 || funcType.NumOut() > 2 {
		return errors.Errorf("Provider function must return 1 or 2 values, but returns %d", funcType.NumOut())
	}

	declaredType := reflect.TypeFor[T]()
	returnType := funcType.Out(0)

	if isTypesNotCompatible(declaredType, returnType) {
		return errors.Errorf("Expected type %s, but function returns %s", funcType.Out(0), returnType.String())
	}

	if funcType.NumOut() == 2 && !funcType.Out(1).Implements(reflect.TypeFor[error]()) {
		return errors.Errorf("Provider function must return an error as second value, but got %s", funcType.Out(1).Kind().String())
	}
	return nil
}

func isTypesNotCompatible(declaredType reflect.Type, returnType reflect.Type) bool {
	if declaredType != returnType {
		if declaredType.Kind() == reflect.Interface && returnType.Kind() == reflect.Ptr {
			return !returnType.AssignableTo(declaredType)
		} else {
			return true
		}
	} else {
		return false
	}
}

func findArgValue(argType reflect.Type, opt *ParameterConfig) (interface{}, error) {
	if argType.Kind() == reflect.Ptr {
		return findPointerType(argType, opt)
	}
	if argType.Kind() == reflect.Interface || argType.Kind() == reflect.Struct {
		return getBeanFromContext(argType)
	} else {
		return findGenricValue(opt)
	}
}

func findPointerType(argType reflect.Type, opt *ParameterConfig) (interface{}, error) {
	err := validateTypeIsBean(argType.Elem())
	if err != nil {
		return nil, err
	}
	bean, err := getBeanFromContext(argType)
	if err != nil {
		if errors.Is(err, ErrNoBeanProvider) {
			return opt.DefaultValue, nil
		}
		return nil, err
	}
	return bean, nil
}

func validateTypeIsBean(argType reflect.Type) error {
	if argType.Kind() == reflect.Struct || argType.Kind() == reflect.Interface {
		return nil
	}
	return errors.Errorf("Expected struct or interface, but got %s", argType.Kind().String())
}

func getBeanFromContext(argType reflect.Type) (interface{}, error) {
	val, err := globalCtx.Get(argType)
	if err != nil {
		return nil, err
	}
	return val, err
}

func findGenricValue(opt *ParameterConfig) (interface{}, error) {
	if opt.ValuePath == "" {
		return reflect.ValueOf(nil), errors.Errorf("Expected non-empty value path")
	}
	value, err := globalCtx.GetGenericValue(opt.ValuePath)
	if err != nil {
		if errors.Is(err, ErrNoValueFound) {
			return opt.DefaultValue, nil
		}
		return nil, err
	}
	return value, nil
}
