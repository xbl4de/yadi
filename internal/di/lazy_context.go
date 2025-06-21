package di

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/internal/types"
	"io"
	"reflect"
	"slices"
	"strings"
)

type LazyContext struct {
	beans       map[reflect.Type]*types.BeanContainer
	providers   map[reflect.Type]*types.BeanProvider
	values      map[string]interface{}
	injectStack []reflect.Type
}

func NewLazyContext(updates []func(ctx types.Context) error) *LazyContext {
	ctx := &LazyContext{
		beans:       make(map[reflect.Type]*types.BeanContainer),
		providers:   make(map[reflect.Type]*types.BeanProvider),
		values:      make(map[string]interface{}),
		injectStack: make([]reflect.Type, 0),
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
	for _, bean := range ctx.beans {
		if !bean.HoldByContext {
			continue
		}
		closeable, ok := bean.Bean.(io.Closer)
		if ok {
			err := closeable.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ctx *LazyContext) addProviders(providers []*types.BeanProvider) {
	for _, provider := range providers {
		ctx.providers[provider.BeanType] = provider
	}
}

func (ctx *LazyContext) Register(provider *types.BeanProvider) error {
	ctx.providers[provider.BeanType] = provider
	return nil
}

func (ctx *LazyContext) Get(p reflect.Type) (types.Bean, error) {
	var err error
	err = ctx.pushInjectStack(p)
	defer ctx.popInjectStack()
	if err != nil {
		return nil, err
	}
	if bean, ok := ctx.beans[p]; ok {
		return bean.Bean, nil
	}
	beanContainer, err := ctx.initBean(p)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to init bean %s", p.String())
	}
	return beanContainer.Bean, nil
}

func (ctx *LazyContext) initBean(p reflect.Type) (*types.BeanContainer, error) {
	var provider *types.BeanProvider
	var ok bool
	var beanContainer *types.BeanContainer
	if provider, ok = ctx.providers[p]; !ok {
		val, err := tryToBuildNewBean(p)
		if err != nil {
			return nil, err
		}
		beanContainer = &types.BeanContainer{
			Bean:          val,
			Type:          p,
			HoldByContext: true,
		}
		return beanContainer, nil
	}

	if provider.UseExistingBean != nil {
		existingBean, err := ctx.Get(provider.UseExistingBean)
		if err != nil {
			return nil, err
		}
		beanContainer = &types.BeanContainer{
			Bean:          existingBean,
			Type:          provider.BeanType,
			HoldByContext: false,
		}
	} else {
		bean, err := provider.Builder(ctx)
		if err != nil {
			return nil, err
		}
		beanContainer = &types.BeanContainer{
			Bean:          bean,
			Type:          provider.BeanType,
			HoldByContext: provider.HoldByContext,
		}
	}
	ctx.beans[provider.BeanType] = beanContainer
	return beanContainer, nil
}

func (ctx *LazyContext) GetGenericValue(path string) (interface{}, error) {
	if val, ok := ctx.values[path]; ok {
		return val, nil
	}
	return nil, types.ErrNoValueFound
}

func (ctx *LazyContext) SetGenericValue(path string, value interface{}) {
	ctx.values[path] = value
}

func (ctx *LazyContext) pushInjectStack(p reflect.Type) error {
	if slices.Contains(ctx.injectStack, p) {
		diStack := ctx.dumpDiStack(p)
		return fmt.Errorf("%w: cannot inject to\n%s", types.ErrCycleDependencies, diStack)
	}
	ctx.injectStack = append(ctx.injectStack, p)
	return nil
}

func (ctx *LazyContext) popInjectStack() {
	if len(ctx.injectStack) == 0 {
		return
	}
	ctx.injectStack = ctx.injectStack[:len(ctx.injectStack)-1]
}

func (ctx *LazyContext) dumpDiStack(toAppend reflect.Type) string {
	builder := strings.Builder{}
	builder.WriteString(ctx.injectStack[0].String() + "\n")
	for _, t := range ctx.injectStack[1:] {
		builder.WriteString("↳ " + t.String() + "\n")
	}
	builder.WriteString("→ " + toAppend.String() + "\n")
	return builder.String()
}
