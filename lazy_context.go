package yadi

import (
	"github.com/pkg/errors"
	"reflect"
)

type LazyContext struct {
	beans     map[reflect.Type]*BeanContainer
	providers map[reflect.Type]*BeanProvider
	values    map[string]interface{}
}

func NewLazyContext(updates []func(ctx Context) error) *LazyContext {
	ctx := &LazyContext{
		beans:     make(map[reflect.Type]*BeanContainer),
		providers: make(map[reflect.Type]*BeanProvider),
		values:    make(map[string]interface{}),
	}
	for _, update := range updates {
		err := update(ctx)
		if err != nil {
			panic(err)
		}
	}
	return ctx
}

func (ctx *LazyContext) Init() {
	// no inits
}

func (ctx *LazyContext) Close() error {
	return nil
}

func (ctx *LazyContext) addProviders(providers []*BeanProvider) {
	for _, provider := range providers {
		ctx.providers[provider.beanType] = provider
	}
}

func (ctx *LazyContext) Register(provider *BeanProvider) error {
	ctx.providers[provider.beanType] = provider
	return nil
}

func (ctx *LazyContext) Get(p reflect.Type) (Bean, error) {
	if bean, ok := ctx.beans[p]; ok {
		return bean, nil
	}
	beanContainer, err := ctx.initBean(p)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to init bean %s", p.String())
	}
	return beanContainer.Bean, nil
}

func (ctx *LazyContext) initBean(p reflect.Type) (*BeanContainer, error) {
	var provider *BeanProvider
	var ok bool
	if provider, ok = ctx.providers[p]; !ok {
		return nil, ErrNoBeanProvider
	}

	var beanContainer *BeanContainer
	if provider.useExistingBean != nil {
		existingBean, err := ctx.Get(provider.useExistingBean)
		if err != nil {
			return nil, err
		}
		beanContainer = &BeanContainer{
			Bean:          existingBean,
			Type:          provider.beanType,
			HoldByContext: false,
		}
	} else {
		bean, err := provider.builder(ctx)
		if err != nil {
			return nil, err
		}
		beanContainer = &BeanContainer{
			Bean:          bean,
			Type:          provider.beanType,
			HoldByContext: provider.holdByContext,
		}
	}
	ctx.beans[provider.beanType] = beanContainer
	return beanContainer, nil
}

func (ctx *LazyContext) GetGenericValue(path string) (interface{}, error) {
	if val, ok := ctx.values[path]; ok {
		return val, nil
	}
	return nil, ErrNoValueFound
}

func (ctx *LazyContext) SetGenericValue(path string, value interface{}) {
	ctx.values[path] = value
}
