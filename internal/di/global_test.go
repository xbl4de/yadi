package di

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/xbl4de/yadi/internal/types"
	"testing"
)

func TestSetBeanProvider_ProviderSet(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return &ServiceE{
			Description: "service E",
		}, nil
	})
	UseLazyContext()

	bean, err := GetBean[*ServiceE]()
	require.NoError(t, err)
	require.NotNil(t, bean)
	require.Equal(t, "service E", bean.Value.Description)
}

func TestSetBeanProvider_ProviderSet_WithHoldByUser(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceClose](func(ctx types.Context) (*ServiceClose, error) {
		return NewServiceClose(), nil
	}, WithHoldByUser())
	UseLazyContext()

	bean, err := GetBean[*ServiceClose]()
	require.NoError(t, err)
	require.NotNil(t, bean)
	require.False(t, bean.Value.Closed)
	err = CloseContext()
	require.NoError(t, err)
	require.False(t, bean.Value.Closed)
}

func TestSetBeanProvider_ProviderSet_WithHoldNotByUser(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceClose](func(ctx types.Context) (*ServiceClose, error) {
		return NewServiceClose(), nil
	})
	UseLazyContext()

	bean, err := GetBean[*ServiceClose]()
	require.NoError(t, err)
	require.NotNil(t, bean)
	require.False(t, bean.Value.Closed)
	err = CloseContext()
	require.NoError(t, err)
	require.True(t, bean.Value.Closed)
}

func TestSetBeanProvider_AutoBuild_ValuesProvided(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	bean, err := GetBean[*ServiceE]()
	require.NoError(t, err)
	require.NotNil(t, bean)
	require.Equal(t, ServiceEDescription, bean.Value.Description)
}

func TestSetBeanProvider_AutoBuild_ValuesNotProvided(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*ServiceE]()
	require.ErrorIs(t, err, types.ErrNoValueFound)
}

func TestGetBean_WithNilContext(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return NewServiceE(ServiceEDescription), nil
	})
	_, err := GetBean[*ServiceE]()
	require.ErrorIs(t, err, types.ErrNilContext)
}

func TestRequireBean_Success(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return NewServiceE(ServiceEDescription), nil
	})
	serviceE := RequireBean[*ServiceE]()
	require.NotNil(t, serviceE)
	require.Equal(t, ServiceEDescription, serviceE.Description)
}

func TestRequireBean_CannotBuild(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.Error(t, r.(error))
		}
	}()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("fail")
	})
	_ = RequireBean[*ServiceE]()
}

func TestRequireBean_NilContext(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.ErrorIs(t, r.(error), types.ErrNilContext)
		}
	}()
	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return NewServiceE("abc"), nil
	})
	_ = RequireBean[*ServiceE]()
}

func TestGetBeanOrDefault_ShouldReturnDefault(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return nil, errors.New("fail")
	})

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})
	require.NotNil(t, bean)
	require.Equal(t, "abcd", bean.Value.Description)
}

func TestGetBeanOrDefault_ShouldReturnRegisteredBean(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx types.Context) (*ServiceE, error) {
		return &ServiceE{Description: "registered"}, nil
	})

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})
	require.NotNil(t, bean)
	require.Equal(t, "registered", bean.Value.Description)
}

func TestGetBeanOrDefault_NilContext(t *testing.T) {
	ResetYadi()

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})
	require.NotNil(t, bean)
	require.Equal(t, "abcd", bean.Value.Description)
}

func TestGetValue_FromPresetDefaults(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	d, err := GetValue[string]("serviceE.description")
	require.NoError(t, err)
	require.Equal(t, ServiceEDescription, d.Value)
}

func TestGetValue_FromProvided(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetValue[string]("serviceE.description", "value")

	d, err := GetValue[string]("serviceE.description")
	require.NoError(t, err)
	require.Equal(t, "value", d.Value)
}

func TestGetValue_GetByWrongType(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetValue[int]("serviceE.description", 10)

	_, err := GetValue[string]("serviceE.description")
	require.Error(t, err)
}

func TestGetValue_NilContext(t *testing.T) {
	ResetYadi()

	SetValue[string]("serviceE.description", "abc")

	_, err := GetValue[string]("serviceE.description")
	require.ErrorIs(t, err, types.ErrNilContext)
}

func TestGetValueOrDefault_ShouldReturnDefault(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	val := GetValueOrDefault[string]("path", "fallback")
	require.NotNil(t, val)
	require.Equal(t, "fallback", val.Value)
}

func TestGetValueOrDefault_ShouldReturnRegistered(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetValue[string]("path", "registered")

	val := GetValueOrDefault[string]("path", "fallback")
	require.NotNil(t, val)
	require.Equal(t, "registered", val.Value)
}

func TestGetValueOrDefault_WrongType_ShouldReturnDefault(t *testing.T) {
	ResetYadi()
	UseLazyContext()

	SetValue[string]("path", "registered")

	val := GetValueOrDefault[int]("path", 10)
	require.NotNil(t, val)
	require.Equal(t, 10, val.Value)
}

func TestGetValueOrDefault_NilContext(t *testing.T) {
	ResetYadi()
	SetValue[int]("path", 10)

	val := GetValueOrDefault[int]("path", 12)
	require.NotNil(t, val)
	require.Equal(t, 12, val.Value)
}

func TestUseLazyContext_UseTwice_ShouldPanic(t *testing.T) {
	ResetYadi()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.ErrorIs(t, r.(error), types.ErrContextAlreadyExists)
		}
	}()

	UseLazyContext()
	UseLazyContext()
}

func TestGetBean_InterfaceType_FromPointerProvider(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[CountInterface](NewServiceF,
		WithDefaultValueAt(0, 22))

	bean, err := GetBean[CountInterface]()
	require.NoError(t, err)
	require.NotNil(t, bean)
	require.Equal(t, 22, bean.Value.GetCount())
}

func TestGetBean_InterfaceType_FromInterfaceProvider(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[CountInterface](func() CountInterface {
		return &ServiceF{
			Count: 11,
		}
	})

	bean, err := GetBean[CountInterface]()
	require.NoError(t, err)
	require.NotNil(t, bean)
	require.Equal(t, 11, bean.Value.GetCount())
}
