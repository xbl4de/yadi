package yadi

import (
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/log"
	"github.com/xbl4de/yadi/types"
	"github.com/xbl4de/yadi/utils"
	"reflect"
)

var globalCtx types.Context

var deferredUpdates []func(ctx types.Context) error

const dummyInt = 42

func getGlobalCtx() types.Context {
	return globalCtx
}

func ensureContext() error {
	if globalCtx == nil {
		return types.ErrNilContext
	}
	return nil
}

func provideDefault[T types.Bean](provider *types.BeanProvider) {
	if globalCtx != nil {
		err := globalCtx.Register(provider)
		if err != nil {
			panic(err)
		}
	} else {
		deferredUpdates = append(deferredUpdates, func(ctx types.Context) error {
			return ctx.Register(provider)
		})
	}
	log.Log("Provided default bean: %s", reflect.TypeFor[T]().String())
}

func clearDeferredUpdates() {
	deferredUpdates = deferredUpdates[:0]
}

func closeContextSoft() error {
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

func applyContext(ctx types.Context) {
	if globalCtx != nil {
		panic(types.ErrContextAlreadyExists)
	}
	globalCtx = ctx
}

func getBeanFromContext(beanType reflect.Type) (types.Bean, error) {
	val, err := getGlobalCtx().Get(beanType)
	if err != nil {
		return nil, err
	}
	return val, err
}

func getGenericValueOrDefault(path string, defaultValue interface{}) (interface{}, error) {
	value, err := getGlobalCtx().GetGenericValue(path)
	if err != nil {
		if errors.Is(err, types.ErrNoValueFound) && defaultValue != nil {
			return defaultValue, nil
		}
		return nil, err
	}
	return value, nil
}

func getGenericValue(path string) (interface{}, error) {
	return getGenericValueOrDefault(path, nil)
}

func getBeanOrDefaultFromContext(beanType reflect.Type, defaultValue types.Bean) (types.Bean, error) {
	err := utils.ValidateTypeIsBean(beanType)
	if err != nil {
		return nil, err
	}
	bean, err := getBeanFromContext(beanType)
	if err != nil {
		if types.ErrNoInjectableProvided(err) && defaultValue != nil {
			return defaultValue, nil
		}
		return nil, err
	}
	return bean, nil
}
