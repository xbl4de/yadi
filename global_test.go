package yadi

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetBeanProvider_ProviderSet(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceE](func(ctx Context) (*ServiceE, error) {
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
	resetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceClose](func(ctx Context) (*ServiceClose, error) {
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
	resetYadi()
	ProvideDefaultValues()

	SetBeanProvider[*ServiceClose](func(ctx Context) (*ServiceClose, error) {
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
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	bean, err := GetBean[*ServiceE]()
	require.NoError(t, err)
	require.NotNil(t, bean)
	require.Equal(t, ServiceEDescription, bean.Value.Description)
}

func TestSetBeanProvider_AutoBuild_ValuesNotProvided(t *testing.T) {
	resetYadi()
	UseLazyContext()

	_, err := GetBean[*ServiceE]()
	require.ErrorIs(t, err, ErrNoValueFound)
}

func TestGetBean_WithNilContext(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, ServiceEDescription))
	_, err := GetBean[*ServiceE]()
	require.ErrorIs(t, err, ErrNilContext)
}

func TestRequireBean_Success(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProviderFunc[*ServiceE](NewServiceE,
		WithDefaultValueAt(0, ServiceEDescription))
	serviceE := RequireBean[*ServiceE]()
	require.NotNil(t, serviceE)
	require.Equal(t, ServiceEDescription, serviceE.Description)
}

func TestRequireBean_CannotBuild(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.Error(t, r.(error))
		}
	}()

	SetBeanProviderFunc[*ServiceE](func() (*ServiceE, error) {
		return nil, errors.New("fail")
	})
	_ = RequireBean[*ServiceE]()
}

func TestRequireBean_NilContext(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.ErrorIs(t, r.(error), ErrNilContext)
		}
	}()

	SetBeanProviderFunc[*ServiceE](NewServiceE)
	_ = RequireBean[*ServiceE]()
}

func TestGetBeanOrDefault_ShouldReturnDefault(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx Context) (*ServiceE, error) {
		return nil, errors.New("fail")
	})

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})
	require.NotNil(t, bean)
	require.Equal(t, "abcd", bean.Value.Description)
}

func TestGetBeanOrDefault_ShouldReturnRegisteredBean(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	SetBeanProvider[*ServiceE](func(ctx Context) (*ServiceE, error) {
		return &ServiceE{Description: "registered"}, nil
	})

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})
	require.NotNil(t, bean)
	require.Equal(t, "registered", bean.Value.Description)
}

func TestGetBeanOrDefault_NilContext(t *testing.T) {
	resetYadi()

	bean := GetBeanOrDefault[*ServiceE](&ServiceE{
		Description: "abcd",
	})
	require.NotNil(t, bean)
	require.Equal(t, "abcd", bean.Value.Description)
}

func TestGetValue_FromPresetDefaults(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	d, err := GetValue[string]("serviceE.description")
	require.NoError(t, err)
	require.Equal(t, ServiceEDescription, d.Value)
}

func TestGetValue_FromProvided(t *testing.T) {
	resetYadi()
	UseLazyContext()

	SetValue[string]("serviceE.description", "value")

	d, err := GetValue[string]("serviceE.description")
	require.NoError(t, err)
	require.Equal(t, "value", d.Value)
}

func TestGetValue_GetByWrongType(t *testing.T) {
	resetYadi()
	UseLazyContext()

	SetValue[int]("serviceE.description", 10)

	_, err := GetValue[string]("serviceE.description")
	require.Error(t, err)
}

func TestGetValue_NilContext(t *testing.T) {
	resetYadi()

	SetValue[string]("serviceE.description", "abc")

	_, err := GetValue[string]("serviceE.description")
	require.ErrorIs(t, err, ErrNilContext)
}

func TestGetValueOrDefault_ShouldReturnDefault(t *testing.T) {
	resetYadi()
	UseLazyContext()

	val := GetValueOrDefault[string]("test.path", "fallback")
	require.NotNil(t, val)
	require.Equal(t, "fallback", val.Value)
}

func TestGetValueOrDefault_ShouldReturnRegistered(t *testing.T) {
	resetYadi()
	UseLazyContext()

	SetValue[string]("test.path", "registered")

	val := GetValueOrDefault[string]("test.path", "fallback")
	require.NotNil(t, val)
	require.Equal(t, "registered", val.Value)
}

func TestGetValueOrDefault_WrongType_ShouldReturnDefault(t *testing.T) {
	resetYadi()
	UseLazyContext()

	SetValue[string]("test.path", "registered")

	val := GetValueOrDefault[int]("test.path", 10)
	require.NotNil(t, val)
	require.Equal(t, 10, val.Value)
}

func TestGetValueOrDefault_NilContext(t *testing.T) {
	resetYadi()
	SetValue[int]("test.path", 10)

	val := GetValueOrDefault[int]("test.path", 12)
	require.NotNil(t, val)
	require.Equal(t, 12, val.Value)
}

func TestUseLazyContext_UseTwice_ShouldPanic(t *testing.T) {
	resetYadi()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.ErrorIs(t, r.(error), ErrContextAlreadyExists)
		}
	}()

	UseLazyContext()
	UseLazyContext()
}

func TestGetBean_InterfaceType_FromPointerProvider(t *testing.T) {
	resetYadi()
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
	resetYadi()
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
