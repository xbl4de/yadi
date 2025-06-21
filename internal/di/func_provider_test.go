package di

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/xbl4de/yadi/internal/types"
	"testing"
)

func TestFuncProviderConfig_Parameter_AtCorrectIndex(t *testing.T) {
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
	require.NotNil(t, p)
	p = cfg.Parameter(1)
	require.NotNil(t, p)
}

func TestFuncProviderConfig_Parameter_AtWrongIndex(t *testing.T) {
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
	require.Equal(t, "", p.ValuePath)
	require.Equal(t, nil, p.DefaultValue)
}

func TestFuncProviderConfig_Parameter_AtNegativeIndex(t *testing.T) {
	opts := []FuncProviderOption{
		func(config *FuncProviderConfig) {
			config.parameters[0] = &ParameterConfig{}
		},
	}
	cfg := NewFuncProviderConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.Error(t, r.(error))
		}
	}()
	cfg.Parameter(-1)
}

func TestSetBeanProviderFunc_WithOption_DefaultValueProvided(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE, WithDefaultValueAt(0, "test"))
	serviceE, err := GetBean[*ServiceE]()
	require.NoError(t, err)
	require.NotNil(t, serviceE)
	require.Equal(t, serviceE.Value.Description, "test")
}

func TestSetBeanProviderFunc_WithOption_DefaultValueNotProvided(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE)
	_, err := GetBean[*ServiceE]()
	require.Error(t, err)
}

func TestSetBeanProviderFunc_WithOption_ValuePathProvided(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithValuePathAt(0, "serviceE.description"))
	SetValue("serviceE.description", "test-description")
	serviceE, err := GetBean[*ServiceE]()
	require.NoError(t, err)
	require.NotNil(t, serviceE)
	require.Equal(t, serviceE.Value.Description, "test-description")
}

func TestSetBeanProviderFunc_WithOption_ValuePathNotProvided(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE)
	_, err := GetBean[*ServiceE]()
	require.Error(t, err)
}

func TestSetBeanProviderFunc_WithErrAsSecondReturn_Success(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](func() (*ServiceE, error) {
		return &ServiceE{}, nil
	})
	serviceE, err := GetBean[*ServiceE]()
	require.NoError(t, err)
	require.NotNil(t, serviceE)
}

func TestSetBeanProviderFunc_WithBeanAsParameter_Success(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceB](NewServiceB,
		WithDefaultValueAt(0, 14))
	serviceB, err := GetBean[*ServiceB]()
	require.NoError(t, err)
	require.NotNil(t, serviceB)
	require.NotNil(t, serviceB.Value.ServiceF)
	require.NotNil(t, serviceB.Value.ServiceH)
	require.Equal(t, 14, serviceB.Value.Age)
}

func TestSetBeanProviderFunc_WithBeanAsParameter_Failure(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceB](NewServiceB,
		WithDefaultValueAt(0, 14))
	SetBeanProvider[*ServiceH](func(ctx types.Context) (*ServiceH, error) {
		return nil, errors.New("test-error")
	})
	_, err := GetBean[*ServiceB]()
	require.Error(t, err)
}

func TestSetBeanProviderFunc_NotFuncProvided(t *testing.T) {
	ResetYadi()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceE](ServiceE{})
	_, err := GetBean[*ServiceE]()
	require.Error(t, err)
}

func TestSetBeanProviderFunc_FuncReturnsNothing(t *testing.T) {
	ResetYadi()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceE](func() {})
	_, err := GetBean[*ServiceE]()
	require.Error(t, err)
}

func TestSetBeanProvider_FuncReturnsTooMuch(t *testing.T) {
	ResetYadi()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceE](func() (*ServiceE, error, bool) {
		return &ServiceE{}, nil, false
	})
	_, err := GetBean[*ServiceE]()
	require.Error(t, err)
}

func TestSetBeanProvider_ReturnTypeMismatch(t *testing.T) {
	ResetYadi()
	UseLazyContext()
	SetBeanProviderFunc[*ServiceE](func() (*ServiceG, error) {
		return &ServiceG{}, nil
	})
	_, err := GetBean[*ServiceE]()
	require.Error(t, err)
}

func TestSetBeanProvider_SecondReturnTypeIsNotError(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](func() (*ServiceE, bool) {
		return &ServiceE{}, false
	})
	_, err := GetBean[*ServiceE]()
	require.Error(t, err)
}
