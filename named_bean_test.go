package yadi

import (
	g "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/types"
	"testing"
)

func TestGetNamedBean_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "abc"),
		WithFuncProviderBeanName("serviceE"))

	serviceE, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceE).ShouldNot(g.BeNil())
}

func TestGetNamedBean_WrongName(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "____"),
		WithFuncProviderBeanName("serviceE"))

	_, err := GetNamedBean[*ServiceE]("bad_name")
	g.Expect(err).Should(g.MatchError(types.ErrNoBeanProvider))
}

func TestGetNamedBean_CannotBuild(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("this is an error")
	}, WithBeanName("serviceE"))

	_, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).Should(g.HaveOccurred())
}

func TestGetNamedBean_CannotBuild_ButAvailableWithoutName(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("this is an error")
	}, WithBeanName("serviceE"))

	bean, err := GetBean[*ServiceE]()
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal(ServiceEDescription))
}

func TestGetNamedBean_NilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()

	_, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).Should(g.MatchError(types.ErrNilContext))
}

func TestGetNamedBean_GetTwice(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "abc"),
		WithFuncProviderBeanName("serviceE"))

	bean, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())

	bean2, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean2).ShouldNot(g.BeNil())

	g.Expect(bean2).Should(g.Equal(bean))
}

func TestGetNamedBean_GetNamedAndUnnamedBean(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "abc"),
		WithFuncProviderBeanName("serviceE"))

	bean, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal("abc"))

	bean2, err := GetBean[*ServiceE]()
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean2).ShouldNot(g.BeNil())
	g.Expect(bean2.Description).Should(g.Equal(ServiceEDescription))

	g.Expect(bean2).ShouldNot(g.Equal(bean))
}

func TestRequireNamedBean_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "abc"),
		WithFuncProviderBeanName("serviceE"))
	serviceE := RequireNamedBean[*ServiceE]("serviceE")
	g.Expect(serviceE).ShouldNot(g.BeNil())
}

func TestRequireNamedBean_WrongName(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "____"),
		WithFuncProviderBeanName("serviceE"))

	g.Expect(func() {
		_ = RequireNamedBean[*ServiceE]("bad_name")
	}).Should(g.PanicWith(g.MatchError(types.ErrNoBeanProvider)))
}

func TestRequireNamedBean_CannotBuild(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("this is an error")
	}, WithBeanName("serviceE"))

	g.Expect(func() {
		_ = RequireNamedBean[*ServiceE]("serviceE")
	}).Should(g.Panic())
}

func TestRequireNamedBean_CannotBuild_ButAvailableWithoutName(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("this is an error")
	}, WithBeanName("serviceE"))

	bean := RequireBean[*ServiceE]()
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal(ServiceEDescription))
}

func TestRequireNamedBean_NilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()

	g.Expect(func() {
		_ = RequireNamedBean[*ServiceE]("serviceE")
	}).Should(g.PanicWith(g.MatchError(types.ErrNilContext)))
}

func TestRequireNamedBean_GetTwice(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "abc"),
		WithFuncProviderBeanName("serviceE"))

	bean, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal("abc"))

	bean2, err := GetNamedBean[*ServiceE]("serviceE")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(bean2).ShouldNot(g.BeNil())
	g.Expect(bean2.Description).Should(g.Equal("abc"))

	g.Expect(bean2).Should(g.Equal(bean))
}

func TestGetNamedBeanOrDefault_ShouldReturnRegistered(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, "abc"),
		WithFuncProviderBeanName("serviceE"))

	bean := GetNamedBeanOrDefault("serviceE", NewServiceE("def"))
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal("abc"))
}

func TestGetNamedBeanOrDefault_ShouldReturnDefault(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	bean := GetNamedBeanOrDefault("serviceE", NewServiceE("def"))
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal("def"))
}

func TestGetNamedBeanOrDefault_CannotBuildRegistered(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("this is an error")
	})

	bean := GetNamedBeanOrDefault("serviceE", NewServiceE("def"))
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal("def"))
}

func TestGetNamedBeanOrDefault_NilContext(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()

	bean := GetNamedBeanOrDefault("serviceE", NewServiceE("def"))
	g.Expect(bean).ShouldNot(g.BeNil())
	g.Expect(bean.Description).Should(g.Equal("def"))
}
