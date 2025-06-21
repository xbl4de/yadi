package di

import (
	"fmt"
	g "github.com/onsi/gomega"
	"github.com/xbl4de/yadi/types"
	"testing"
)

type A struct {
	C *C
}

type B struct {
	A *A
}

type C struct {
	B *B
}

type NewA struct {
	C *LazyC
}

type NewB struct {
	A *NewA
}

type LazyC struct {
	B types.LazyBean[*NewB]
}

func TestCycleDependencies_BeanA(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*A]()

	g.Expect(err).Should(g.MatchError(types.ErrCycleDependencies))
	fmt.Println(err)
}

func TestCycleDependencies_BeanB(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*B]()

	g.Expect(err).Should(g.MatchError(types.ErrCycleDependencies))
}

func TestCycleDependencies_BeanC(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*C]()

	g.Expect(err).Should(g.MatchError(types.ErrCycleDependencies))
}

func TestCycleDependencies_BeanLazy(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProvider[*LazyC](func(ctx types.Context) (*LazyC, error) {
		return &LazyC{
			B: NewLazyBean[*NewB](),
		}, nil
	})

	c, err := GetBean[*LazyC]()

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(c).ShouldNot(g.BeNil())
	g.Expect(c.Value.B().A).ShouldNot(g.BeNil())
	g.Expect(c.Value.B().A.C).ShouldNot(g.BeNil())
}
