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
	box := buildValueBox[T](outs[0])

	if len(outs) == 1 {
		return box, nil
	} else {
		return box, castToErr(outs[1])
	}
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

func findArgValue(argType reflect.Type, opt *ParameterConfig) (interface{}, error) {
	if argType.Kind() == reflect.Ptr ||
		argType.Kind() == reflect.Interface ||
		argType.Kind() == reflect.Struct {
		return getBeanOrDefault(argType, opt.DefaultValue)
	}
	return getGenericValueOrDefault(opt.ValuePath, opt.DefaultValue)
}
