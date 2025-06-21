package di

import (
	"fmt"
	g "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/types"
	"reflect"
	"testing"
)

func requireType[T types.Bean](val interface{}) T {
	casted, ok := val.(T)
	g.Expect(ok).Should(g.BeTrue())
	return casted
}

func TestAutoBeanPtrGeneration_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	serviceABean, err := tryToBuildNewBean(reflect.TypeFor[*ServiceA]())
	serviceA := requireType[*ServiceA](serviceABean)

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceA).ShouldNot(g.BeNil())
	g.Expect(serviceA.ServiceF).ShouldNot(g.BeNil())
	g.Expect(serviceA.ServiceF).ShouldNot(g.BeNil())
	g.Expect(serviceA.Name).Should(g.Equal(ServiceAName))
}

func TestAutoBeanPtrGeneration_WhenCannotBuildInnerBean(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	var errServiceF = errors.New("errServiceF")
	SetBeanProvider(func(ctx types.Context) (*ServiceF, error) {
		return nil, errServiceF
	})
	ProvideDefaultValues()
	UseLazyContext()

	_, err := tryToBuildNewBean(reflect.TypeFor[*ServiceA]())
	g.Expect(err).Should(g.MatchError(errServiceF))
}

func TestAutoBeanGeneration_Success(t *testing.T) {
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	serviceABean, err := tryToBuildNewBean(reflect.TypeFor[ServiceA]())
	serviceA := requireType[ServiceA](serviceABean)

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceA).ShouldNot(g.BeNil())
	g.Expect(serviceA.ServiceF).ShouldNot(g.BeNil())
	g.Expect(serviceA.ServiceF).ShouldNot(g.BeNil())
	g.Expect(serviceA.Name).Should(g.Equal(ServiceAName))
}

func TestInjectToNotBean_ShouldFail(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	abc := ""
	err := Inject(&abc)

	g.Expect(err).Should(g.MatchError(types.ErrInjectNotSupported))
}

func TestInjectToNotPtr_ShouldFail(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	serviceA := ServiceA{}
	err := Inject(serviceA)

	g.Expect(err).Should(g.MatchError(types.ErrInjectNotSupported))
}

func TestInjectToPtr_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	serviceA := ServiceA{}
	err := Inject(&serviceA)

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceA).ShouldNot(g.BeNil())
	g.Expect(serviceA.ServiceF).ShouldNot(g.BeNil())
	g.Expect(serviceA.ServiceF).ShouldNot(g.BeNil())
	g.Expect(serviceA.Name).Should(g.Equal(ServiceAName))
}

func TestInjectToIgnoreTag_ShouldSkipInjection(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	serviceIgnore := ServiceIgnore{}
	err := Inject(&serviceIgnore)

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(serviceIgnore).ShouldNot(g.BeNil())
	g.Expect(serviceIgnore.ServiceA).Should(g.BeNil())
}

func TestTryToBuildNonBeanTypes_ShouldFail(t *testing.T) {
	g.RegisterTestingT(t)
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
			g.Expect(err).Should(g.MatchError(types.ErrNonBeanType))
		})
	}
}

func TestInject_BadTag(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	serviceBagTag := ServiceBagTag{}
	err := Inject(&serviceBagTag)

	g.Expect(err).Should(g.MatchError(types.ErrParseTag))
}

func TestInject_NoValueProvided(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	serviceA := ServiceA{}
	err := Inject(&serviceA)

	g.Expect(err).Should(g.MatchError(types.ErrNoValueFound))
}

func TestInjectLazyBean_Success(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()

	serviceA := InjectLazyBean[*ServiceA]()

	g.Expect(serviceA).ShouldNot(g.BeNil())

	val := serviceA()

	g.Expect(val).ShouldNot(g.BeNil())
	g.Expect(val.ServiceF).ShouldNot(g.BeNil())
	g.Expect(val.ServiceF).ShouldNot(g.BeNil())
	g.Expect(val.Name).Should(g.Equal(ServiceAName))
}

func TestInjectLazyBean_CannotBuildBean(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	ProvideDefaultValues()
	UseLazyContext()
	err := errors.New("error")

	SetBeanProvider[*ServiceA](func(ctx types.Context) (*ServiceA, error) {
		return nil, err
	})

	serviceA := InjectLazyBean[*ServiceA]()

	g.Expect(serviceA).ShouldNot(g.BeNil())
	g.Expect(func() {
		_ = serviceA()
	}).Should(g.PanicWith(g.MatchError(err)))
}
