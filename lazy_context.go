package yadi

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/types"
	"io"
	"reflect"
	"slices"
	"strings"
)

type BeanKey struct {
	Type reflect.Type
	Name string
}

func NewBeanKey(typ reflect.Type, name string) BeanKey {
	return BeanKey{
		Type: typ,
		Name: name,
	}
}

type LazyContext struct {
	beans       map[BeanKey]*types.BeanContainer
	providers   map[BeanKey]*types.BeanProvider
	values      map[string]interface{}
	injectStack []BeanKey
}

func NewLazyContext(updates []func(ctx types.Context) error) *LazyContext {
	ctx := &LazyContext{
		beans:       make(map[BeanKey]*types.BeanContainer),
		providers:   make(map[BeanKey]*types.BeanProvider),
		values:      make(map[string]interface{}),
		injectStack: make([]BeanKey, 0),
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
		key := keyFromProvider(provider)
		ctx.providers[key] = provider
	}
}

func keyFromProvider(provider *types.BeanProvider) BeanKey {
	return NewBeanKey(provider.BeanType, provider.BeanName)
}

func (ctx *LazyContext) Register(provider *types.BeanProvider) error {
	key := keyFromProvider(provider)
	ctx.providers[key] = provider
	return nil
}

func (ctx *LazyContext) Get(typ reflect.Type) (types.Bean, error) {
	return ctx.get(NewBeanKey(typ, ""), true)
}
func (ctx *LazyContext) GetNamed(typ reflect.Type, beanName string) (types.Bean, error) {
	return ctx.get(NewBeanKey(typ, beanName), false)
}

func (ctx *LazyContext) get(key BeanKey, buildIfNotFound bool) (types.Bean, error) {
	var err error
	err = ctx.pushInjectStack(key)
	defer ctx.popInjectStack()
	if err != nil {
		return nil, err
	}
	if bean, ok := ctx.beans[key]; ok {
		return bean.Bean, nil
	}
	beanContainer, err := ctx.initBean(key, buildIfNotFound)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to init bean %s[%s]", key.Name, key.Type.String())
	}
	return beanContainer.Bean, nil
}

func (ctx *LazyContext) pushInjectStack(key BeanKey) error {
	if slices.Contains(ctx.injectStack, key) {
		diStack := ctx.dumpDiStack(key)
		return fmt.Errorf("%w: cannot inject to\n%s", types.ErrCycleDependencies, diStack)
	}
	ctx.injectStack = append(ctx.injectStack, key)
	return nil
}

func (ctx *LazyContext) popInjectStack() {
	if len(ctx.injectStack) == 0 {
		return
	}
	ctx.injectStack = ctx.injectStack[:len(ctx.injectStack)-1]
}

func (ctx *LazyContext) dumpDiStack(toAppend BeanKey) string {
	builder := strings.Builder{}
	firstEl := ctx.injectStack[0]
	builder.WriteString(fmt.Sprintf("%s[%s]", firstEl.Name, firstEl.Type.String()))
	for _, el := range ctx.injectStack[1:] {
		builder.WriteString(fmt.Sprintf("↳  %s[%s]", el.Name, el.Type.String()))
	}
	builder.WriteString(fmt.Sprintf("→ %s[%s]", toAppend.Name, toAppend.Type.String()))
	return builder.String()
}

func (ctx *LazyContext) initBean(key BeanKey, shouldTryBuildNewBean bool) (*types.BeanContainer, error) {
	var provider *types.BeanProvider
	var ok bool
	var beanContainer *types.BeanContainer
	if provider, ok = ctx.providers[key]; !ok {
		if !shouldTryBuildNewBean {
			return nil, types.ErrNoBeanProvider
		}
		val, err := tryToBuildNewBean(key.Type)
		if err != nil {
			return nil, err
		}
		beanContainer = types.NewBeanContainerHoldByContext(val, key.Name, key.Type)
		return beanContainer, nil
	}

	if provider.UseExistingBean != nil {
		existingBean, err := ctx.get(NewBeanKey(provider.UseExistingBean, key.Name), false)
		if err != nil {
			return nil, err
		}
		beanContainer = types.NewBeanContainerHoldByUser(existingBean, key.Name, key.Type)
	} else {
		bean, err := provider.Builder(ctx)
		if err != nil {
			return nil, err
		}
		beanContainer = types.NewBeanContainer(bean, key.Name, key.Type, provider.HoldByContext)
	}
	ctx.beans[key] = beanContainer
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
