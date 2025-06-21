package types

import (
	g "github.com/onsi/gomega"
	"testing"
)

func TestParseTag_EmptyTag(t *testing.T) {
	g.RegisterTestingT(t)

	tag, err := ParseTag("")

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(tag).ShouldNot(g.BeNil())
	g.Expect(tag.BeanName).Should(g.BeEmpty())
	g.Expect(tag.ValuePath).Should(g.BeEmpty())
	g.Expect(tag.Ignore).Should(g.BeFalse())
}

func TestParseTag_FullTag(t *testing.T) {
	g.RegisterTestingT(t)

	tag, err := ParseTag("ignore;path=abc;beanName=def")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(tag).ShouldNot(g.BeNil())
	g.Expect(tag.ValuePath).Should(g.Equal("abc"))
	g.Expect(tag.BeanName).Should(g.Equal("def"))
	g.Expect(tag.Ignore).Should(g.BeTrue())
}

func TestParseTag_IgnoreTag(t *testing.T) {
	g.RegisterTestingT(t)

	tag, err := ParseTag("ignore")
	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(tag).ShouldNot(g.BeNil())
	g.Expect(tag.ValuePath).Should(g.BeEmpty())
	g.Expect(tag.BeanName).Should(g.BeEmpty())
}

func TestParseTag_BeanName(t *testing.T) {
	g.RegisterTestingT(t)
	tag, err := ParseTag("beanName=beanA")

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(tag).ShouldNot(g.BeNil())
	g.Expect(tag.BeanName).Should(g.Equal("beanA"))
	g.Expect(tag.ValuePath).Should(g.BeEmpty())
	g.Expect(tag.Ignore).Should(g.BeFalse())
}

func TestParseTag_ValuePath(t *testing.T) {
	g.RegisterTestingT(t)
	tag, err := ParseTag("path=abc")

	g.Expect(err).ShouldNot(g.HaveOccurred())
	g.Expect(tag).ShouldNot(g.BeNil())
	g.Expect(tag.ValuePath).Should(g.Equal("abc"))
	g.Expect(tag.Ignore).Should(g.BeFalse())
	g.Expect(tag.BeanName).Should(g.BeEmpty())
}

func TestParseTag_EmptyValuePath(t *testing.T) {
	g.RegisterTestingT(t)
	_, err := ParseTag("path=")

	g.Expect(err).Should(g.MatchError(ErrParseTag))
}

func TestParseTag_EmptyBeanName(t *testing.T) {
	g.RegisterTestingT(t)
	_, err := ParseTag("beanName=")

	g.Expect(err).Should(g.MatchError(ErrParseTag))
}

func TestParseTag_UnknownTag(t *testing.T) {
	g.RegisterTestingT(t)
	_, err := ParseTag("unknown")

	g.Expect(err).Should(g.HaveOccurred())
}

func TestParseTag_ExtraValues(t *testing.T) {
	g.RegisterTestingT(t)
	_, err := ParseTag("beanName=beanA=beanB")

	g.Expect(err).Should(g.MatchError(ErrParseTag))
}
