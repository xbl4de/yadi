package di

import (
	g "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/internal/types"
	"testing"
)

func TestSetBeanProvider_ProviderSet(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return &ServiceE{
			Description: "service E",
		}, nil
	})
	UseLazyContext()

	bean, err := GetBean[*ServiceE]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.Description).Should(g.Equal("service E"))
}

func TestSetBeanProvider_ProviderSet_WithHoldByUser(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceClose](func(ctx types.Context) (*ServiceClose, error) {
		return NewServiceClose(), nil
	}, WithHoldByUser())
	UseLazyContext()

	bean, err := GetBean[*ServiceClose]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.Closed).Should(g.BeFalse())

	err = CloseContext()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean.Value.Closed).Should(g.BeFalse())
}

func TestSetBeanProvider_ProviderSet_WithHoldNotByUser(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceClose](func(ctx types.Context) (*ServiceClose, error) {
		return NewServiceClose(), nil
	})
	UseLazyContext()

	bean, err := GetBean[*ServiceClose]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.Closed).Should(g.BeFalse())

	err = CloseContext()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean.Value.Closed).Should(g.BeTrue())
}

func TestSetBeanProvider_AutoBuild_ValuesProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	bean, err := GetBean[*ServiceE]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.Description).Should(g.Equal(ServiceEDescription))
}

func TestSetBeanProvider_AutoBuild_ValuesNotProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.MatchError(types.ErrNoValueFound))
}

func TestGetBean_WithNilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return NewServiceE(ServiceEDescription), nil
	})
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.MatchError(types.ErrNilContext))
}

func TestRequireBean_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return NewServiceE(ServiceEDescription), nil
	})
	serviceE := RequireBean[*ServiceE]()

	g.Expect(serviceE).ShouldNot(g.BeNil())
	g.Expect(serviceE.Description).Should(g.Equal(ServiceEDescription))
}

func TestRequireBean_CannotBuild(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("fail")
	})

	g.Expect(func() {
		RequireBean[*ServiceE]()
	}).Should(g.Panic())
}

func TestRequireBean_NilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return NewServiceE("abc"), nil
	})

	g.Expect(func() {
		RequireBean[*ServiceE]()
	}).Should(g.PanicWith(g.MatchError(types.ErrNilContext)))
}

func TestGetBeanOrDefault_ShouldReturnDefault(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("fail")
	})

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})

	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.Description).Should(g.Equal("abcd"))
}

func TestGetBeanOrDefault_ShouldReturnRegisteredBean(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return &ServiceE{Description: "registered"}, nil
	})

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})

	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.Description).Should(g.Equal("registered"))
}

func TestGetBeanOrDefault_NilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})

	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.Description).Should(g.Equal("abcd"))
}

func TestGetValue_FromPresetDefaults(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	d, err := GetValue[string]("serviceE.description")

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(d.Value)
}

func TestGetValue_FromProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetValue[string]("serviceE.description", "value")

	d, err := GetValue[string]("serviceE.description")

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(d.Value).Should(g.Equal("value"))
}

func TestGetValue_GetByWrongType(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetValue[int]("serviceE.description", 10)

	_, err := GetValue[string]("serviceE.description")

	g.Expect(err).Should(g.HaveOccurred())
}

func TestGetValue_NilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()

	SetValue[string]("serviceE.description", "abc")

	_, err := GetValue[string]("serviceE.description")

	g.Expect(err).Should(g.MatchError(types.ErrNilContext))
}

func TestGetValueOrDefault_ShouldReturnDefault(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	val := GetValueOrDefault[string]("path", "fallback")

	g.Expect(val).ShouldNot(g.BeNil())
	g.Expect(val.Value).Should(g.Equal("fallback"))
}

func TestGetValueOrDefault_ShouldReturnRegistered(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetValue[string]("path", "registered")

	val := GetValueOrDefault[string]("path", "fallback")

	g.Expect(val).ShouldNot(g.BeNil())
	g.Expect(val.Value).Should(g.Equal("registered"))
}

func TestGetValueOrDefault_WrongType_ShouldReturnDefault(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetValue[string]("path", "registered")

	val := GetValueOrDefault[int]("path", 10)

	g.Expect(val).ShouldNot(g.BeNil())
	g.Expect(val.Value).Should(g.Equal(10))
}

func TestGetValueOrDefault_NilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	SetValue[int]("path", 10)

	val := GetValueOrDefault[int]("path", 12)

	g.Expect(val).ShouldNot(g.BeNil())
	g.Expect(val.Value).Should(g.Equal(12))
}

func TestUseLazyContext_UseTwice_ShouldPanic(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()

	g.Expect(func() {
		UseLazyContext()
		UseLazyContext()
	}).Should(g.Panic())
}

func TestGetBean_InterfaceType_FromPointerProvider(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[CountInterface](NewServiceF,
		WithDefaultValueAt(0, 22))

	bean, err := GetBean[CountInterface]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.GetCount()).Should(g.Equal(22))
}

func TestGetBean_InterfaceType_FromInterfaceProvider(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[CountInterface](func() CountInterface {
		return &ServiceF{
			Count: 11,
		}
	})

	bean, err := GetBean[CountInterface]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Value.GetCount()).Should(g.Equal(11))
}
