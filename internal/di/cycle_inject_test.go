package di

import (
	"fmt"
	g "github.com/onsi/gomega"
	types2 "github.com/xbl4de/yadi/types"
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
	B types2.LazyBean[*NewB]
}

func TestCycleDependencies_BeanA(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*A]()

	g.Expect(err).Should(g.MatchError(types2.ErrCycleDependencies))
	fmt.Println(err)
}

func TestCycleDependencies_BeanB(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*B]()

	g.Expect(err).Should(g.MatchError(types2.ErrCycleDependencies))
}

func TestCycleDependencies_BeanC(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	_, err := GetBean[*C]()

	g.Expect(err).Should(g.MatchError(types2.ErrCycleDependencies))
}

func TestCycleDependencies_BeanLazy(t *testing.T) {
	g.RegisterTestingT(t)
	ResetYadi()
	UseLazyContext()

	SetBeanProvider[*LazyC](func(ctx types2.Context) (*LazyC, error) {
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
