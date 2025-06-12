package yadi

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func requireType[T Bean](t *testing.T, val interface{}) T {
	casted, ok := val.(T)
	require.True(t, ok)
	return casted
}

func TestAutoBeanPtrGeneration_Success(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	serviceABean, err := tryToBuildNewBean(reflect.TypeFor[*ServiceA]())
	serviceA := requireType[*ServiceA](t, serviceABean)
	require.NoError(t, err)
	require.NotNil(t, serviceA)
	require.NotNil(t, serviceA.ServiceE)
	require.NotNil(t, serviceA.ServiceF)
	require.Equal(t, ServiceAName, serviceA.Name)
}

func TestAutoBeanPtrGeneration_WhenCannotBuildInnerBean(t *testing.T) {
	resetYadi()
	var errServiceF = errors.New("errServiceF")
	SetBeanProvider(func(ctx Context) (*ServiceF, error) {
		return nil, errServiceF
	})
	ProvideDefaultValues()
	UseLazyContext()
	serviceA, err := tryToBuildNewBean(reflect.TypeFor[*ServiceA]())
	require.ErrorIs(t, err, errServiceF)
	require.Nil(t, serviceA)
}

func TestAutoBeanGeneration_Success(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	serviceABean, err := tryToBuildNewBean(reflect.TypeFor[ServiceA]())
	serviceA := requireType[ServiceA](t, serviceABean)
	require.NoError(t, err)
	require.NotNil(t, serviceA.ServiceE)
	require.NotNil(t, serviceA.ServiceF)
	require.Equal(t, ServiceAName, serviceA.Name)
}

func TestInjectToNotBean_ShouldFail(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	abc := ""
	err := Inject(&abc)
	require.ErrorIs(t, err, ErrInjectNotSupported)
}

func TestInjectToNotPtr_ShouldFail(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	serviceA := ServiceA{}
	err := Inject(serviceA)
	require.ErrorIs(t, err, ErrInjectNotSupported)
}

func TestInjectToPtr_Success(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	serviceA := ServiceA{}
	err := Inject(&serviceA)
	require.NoError(t, err)
	require.NotNil(t, serviceA.ServiceE)
	require.NotNil(t, serviceA.ServiceF)
	require.Equal(t, ServiceAName, serviceA.Name)
}

func TestInjectToIgnoreTag_ShouldSkipInjection(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	serviceIgnore := ServiceIgnore{}
	err := Inject(&serviceIgnore)
	require.NoError(t, err)
	require.Nil(t, serviceIgnore.ServiceA)
}

func TestTryToBuildNonBeanTypes_ShouldFail(t *testing.T) {
	tests := []struct {
		T reflect.Type
	}{
		{reflect.TypeFor[int]()},
		{reflect.TypeFor[string]()},
		{reflect.TypeFor[bool]()},
		{reflect.TypeFor[float32]()},
		{reflect.TypeFor[float64]()},
		{reflect.TypeFor[[]string]()},
		{reflect.TypeFor[complex128]()},
		{reflect.TypeFor[map[string]string]()},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_ShouldFail", test.T.String()), func(t *testing.T) {
			_, err := tryToBuildNewBean(test.T)
			require.ErrorIs(t, err, ErrNonBeanType)
		})
	}
}

func TestInject_BadTag(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	serviceBagTag := ServiceBagTag{}
	err := Inject(&serviceBagTag)
	require.ErrorIs(t, err, ErrParseTag)
}

func TestInject_NoValueProvided(t *testing.T) {
	resetYadi()
	UseLazyContext()
	serviceA := ServiceA{}
	err := Inject(&serviceA)
	require.ErrorIs(t, err, ErrNoValueFound)
}

func TestInjectLazyBean_Success(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	serviceA := InjectLazyBean[*ServiceA]()
	require.NotNil(t, serviceA)
	val := serviceA()
	require.NotNil(t, val)
	require.Equal(t, ServiceAName, val.Name)
	require.NotNil(t, val.ServiceE)
	require.NotNil(t, val.ServiceF)
}

func TestInjectLazyBean_CannotBuildBean(t *testing.T) {
	resetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	err := errors.New("error")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			require.ErrorIs(t, r.(error), err)
		}
	}()

	SetBeanProvider[*ServiceA](func(ctx Context) (*ServiceA, error) {
		return nil, err
	})

	serviceA := InjectLazyBean[*ServiceA]()
	require.NotNil(t, serviceA)
	_ = serviceA()
}
