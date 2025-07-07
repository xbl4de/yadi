package example

import (
	"fmt"
	g "github.com/onsi/gomega"
	"github.com/xbl4de/yadi"
	"github.com/xbl4de/yadi/types"
	"testing"
)

type ExampleServiceA struct {
	Timeout     int    `yadi:"path=serviceA.timeout"`
	Credentials string `yadi:"path=serviceA.credentials"`
}

type ExampleServiceB struct {
	Timeout  int `yadi:"path=serviceB.timeout"`
	ServiceA *ExampleServiceA
}

func TestAutoBeanProvideExample(t *testing.T) {
	g.RegisterTestingT(t)
	_ = yadi.CloseContext()

	yadi.UseLazyContext()
	yadi.SetValue("serviceA.timeout", 10)
	yadi.SetValue("serviceB.timeout", 20)
	yadi.SetValue("serviceA.credentials", `{"user":"password"}`)

	b, err := yadi.GetBean[*ExampleServiceB]()
	g.Expect(err).ShouldNot(g.HaveOccurred())

	fmt.Printf("Service B: %+v\n", b)
	fmt.Printf("Service A: %+v\n", b.ServiceA)
}

func TestBeanProviderExample(t *testing.T) {
	g.RegisterTestingT(t)
	_ = yadi.CloseContext()

	yadi.UseLazyContext()
	yadi.SetBeanProvider[*ExampleServiceB](func(ctx types.Context) (*ExampleServiceB, error) {
		return &ExampleServiceB{
			ServiceA: &ExampleServiceA{
				Timeout:     15,
				Credentials: `{"user1":"password1"}`,
			},
			Timeout: 25,
		}, nil
	})

	b, err := yadi.GetBean[*ExampleServiceB]()
	g.Expect(err).ShouldNot(g.HaveOccurred())

	fmt.Printf("Service B: %+v\n", b)
	fmt.Printf("Service A: %+v\n", b.ServiceA)
}

func TestFuncBeanProviderExample(t *testing.T) {
	g.RegisterTestingT(t)
	_ = yadi.CloseContext()

	yadi.UseLazyContext()
	provider := func() *ExampleServiceB {
		return &ExampleServiceB{
			ServiceA: &ExampleServiceA{
				Timeout:     5,
				Credentials: `{"user2":"password2"}`,
			},
			Timeout: 35,
		}
	}

	yadi.SetBeanProviderFunc[*ExampleServiceB](provider)

	b, err := yadi.GetBean[*ExampleServiceB]()
	g.Expect(err).ShouldNot(g.HaveOccurred())
	fmt.Printf("Service B: %+v\n", b)
	fmt.Printf("Service A: %+v\n", b.ServiceA)
}

type ExampleServiceC struct {
	unexportedField string `yadi:"path=serviceC.unexportedField"`
}

func (c *ExampleServiceC) SetUnexportedField(unexportedField string) {
	c.unexportedField = unexportedField
}

func TestInjectViaSetterExample(t *testing.T) {
	g.RegisterTestingT(t)
	_ = yadi.CloseContext()

	yadi.UseLazyContext()
	yadi.SetValue("serviceC.unexportedField", "unexportedField")

	bean, err := yadi.GetBean[*ExampleServiceC]()
	g.Expect(err).ShouldNot(g.HaveOccurred())
	fmt.Printf("Service C: %+v\n", bean)
}
