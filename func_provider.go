package yadi

import (
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/types"
	"github.com/xbl4de/yadi/utils"
	"reflect"
)

type ParameterConfig struct {
	ValuePath    string
	DefaultValue interface{}
}

type FuncProviderConfig struct {
	beanName   string
	parameters map[int]*ParameterConfig
}

type FuncProviderOption func(*FuncProviderConfig)

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

func extractBeanName(opts []FuncProviderOption) string {
	beanName := ""
	if len(opts) > 0 {
		fakeCfg := NewFuncProviderConfig()
		for _, opt := range opts {
			opt(fakeCfg)
		}
		beanName = fakeCfg.beanName
	}
	return beanName
}

func providerFromFuncE[T types.Bean](function interface{}, cfg *FuncProviderConfig) (T, error) {
	funcValue := reflect.ValueOf(function)
	funcType := reflect.TypeOf(function)
	var zeroValue T
	err := validateProviderFunc[T](funcValue, funcType)
	if err != nil {
		return zeroValue, err
	}

	args, err := buildArgs(funcType, cfg)
	if err != nil {
		return zeroValue, err
	}

	outs := funcValue.Call(args)
	bean := utils.ConvertToBean(outs[0])

	casted, ok := bean.(T)
	if !ok {
		return zeroValue, errors.Errorf("invalid bean type: %T", bean)
	}

	if len(outs) == 1 {
		return casted, nil
	} else {
		return casted, utils.CastToErr(outs[1])
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

	if utils.IsTypesNotCompatible(declaredType, returnType) {
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
		return getBeanOrDefaultFromContext(argType, opt.DefaultValue)
	}
	return getGenericValueOrDefault(opt.ValuePath, opt.DefaultValue)
}
