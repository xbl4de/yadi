package di

import (
	g "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/types"
	"testing"
)

func TestFuncProviderConfig_Parameter_AtCorrectIndex(t *testing.T) {
	g.RegisterTestingT(t)
	opts := []FuncProviderOption{
		func(config *FuncProviderConfig) {
			config.parameters[0] = &ParameterConfig{}
			config.parameters[1] = &ParameterConfig{}
		},
	}
	cfg := NewFuncProviderConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	p := cfg.Parameter(0)
	g.Expect(p).ShouldNot(g.BeNil())

	p = cfg.Parameter(1)
	g.Expect(p).ShouldNot(g.BeNil())
}

func TestFuncProviderConfig_Parameter_AtWrongIndex(t *testing.T) {
	g.RegisterTestingT(t)
	opts := []FuncProviderOption{
		func(config *FuncProviderConfig) {
			config.parameters[0] = &ParameterConfig{}
		},
	}
	cfg := NewFuncProviderConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	p := cfg.Parameter(1)

	g.Expect(p.ValuePath).Should(g.BeEmpty())
	g.Expect(p.DefaultValue).Should(g.BeNil())
}

func TestFuncProviderConfig_Parameter_AtNegativeIndex(t *testing.T) {
	g.RegisterTestingT(t)
	opts := []FuncProviderOption{
		func(config *FuncProviderConfig) {
			config.parameters[0] = &ParameterConfig{}
		},
	}
	cfg := NewFuncProviderConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	g.Expect(func() {
		cfg.Parameter(-1)
	}).Should(g.Panic())
}

func TestSetBeanProviderFunc_WithOption_DefaultValueProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE, WithDefaultValueAt(0, "test"))
	serviceE, err := GetBean[*ServiceE]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceE).ShouldNot(g.BeNil())
	g.Expect(serviceE.Value.Description).Should(g.Equal("test"))
}

func TestSetBeanProviderFunc_WithOption_DefaultValueNotProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE)
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.HaveOccurred())
}

func TestSetBeanProviderFunc_WithOption_ValuePathProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithValuePathAt(0, "serviceE.description"))
	SetValue("serviceE.description", "test-description")
	serviceE, err := GetBean[*ServiceE]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceE).ShouldNot(g.BeNil())
	g.Expect(serviceE.Value.Description).Should(g.Equal("test-description"))
}

func TestSetBeanProviderFunc_WithOption_ValuePathNotProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE)
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.HaveOccurred())
}

func TestSetBeanProviderFunc_WithErrAsSecondReturn_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](func() (*ServiceE, error) {
		return &ServiceE{}, nil
	})
	serviceE, err := GetBean[*ServiceE]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceE).ShouldNot(g.BeNil())
}

func TestSetBeanProviderFunc_WithBeanAsParameter_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceB](NewServiceB,
		WithDefaultValueAt(0, 14))
	serviceB, err := GetBean[*ServiceB]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceB).ShouldNot(g.BeNil())
	g.Expect(serviceB.Value.ServiceF).ShouldNot(g.BeNil())
	g.Expect(serviceB.Value.ServiceH).ShouldNot(g.BeNil())
	g.Expect(serviceB.Value.Age).Should(g.Equal(14))
}

func TestSetBeanProviderFunc_WithBeanAsParameter_Failure(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceB](NewServiceB,
		WithDefaultValueAt(0, 14))
	SetBeanProvider[*ServiceH](func(ctx types.Context) (*ServiceH, error) {
		return nil, errors.New("test-error")
	})
	_, err := GetBean[*ServiceB]()

	g.Expect(err).Should(g.HaveOccurred())
}

func TestSetBeanProviderFunc_NotFuncProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](ServiceE{})
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.HaveOccurred())
}

func TestSetBeanProviderFunc_FuncReturnsNothing(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceE](func() {})
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.HaveOccurred())
}

func TestSetBeanProvider_FuncReturnsTooMuch(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceE](func() (*ServiceE, error, bool) {
		return &ServiceE{}, nil, false
	})
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.HaveOccurred())
}

func TestSetBeanProvider_ReturnTypeMismatch(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](func() (*ServiceG, error) {
		return &ServiceG{}, nil
	})
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.HaveOccurred())
}

func TestSetBeanProvider_SecondReturnTypeIsNotError(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](func() (*ServiceE, bool) {
		return &ServiceE{}, false
	})
	_, err := GetBean[*ServiceE]()

	g.Expect(err).Should(g.HaveOccurred())
}
