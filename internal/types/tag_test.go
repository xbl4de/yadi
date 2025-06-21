package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseTag_EmptyTag(t *testing.T) {
	tag, err := ParseTag("")
	require.NoError(t, err)
	require.Equal(t, "", tag.BeanName)
	require.Equal(t, false, tag.Ignore)
	require.Equal(t, "", tag.ValuePath)
}

func TestParseTag_FullTag(t *testing.T) {
	tag, err := ParseTag("ignore;path=abc;beanName=def")
	require.NoError(t, err)
	require.Equal(t, "abc", tag.ValuePath)
	require.Equal(t, "def", tag.BeanName)
	require.Equal(t, true, tag.Ignore)
}

func TestParseTag_IgnoreTag(t *testing.T) {
	tag, err := ParseTag("ignore")
	require.NoError(t, err)
	require.Equal(t, "", tag.BeanName)
	require.Equal(t, "", tag.ValuePath)
	require.Equal(t, true, tag.Ignore)
}

func TestParseTag_BeanName(t *testing.T) {
	tag, err := ParseTag("beanName=beanA")
	require.NoError(t, err)
	require.Equal(t, "beanA", tag.BeanName)
	require.Equal(t, false, tag.Ignore)
	require.Equal(t, "", tag.ValuePath)
}

func TestParseTag_ValuePath(t *testing.T) {
	tag, err := ParseTag("path=abc")
	require.NoError(t, err)
	require.Equal(t, "abc", tag.ValuePath)
	require.Equal(t, false, tag.Ignore)
	require.Equal(t, "", tag.BeanName)
}

func TestParseTag_EmptyValuePath(t *testing.T) {
	_, err := ParseTag("path=")
	require.ErrorIs(t, err, ErrParseTag)
}

func TestParseTag_EmptyBeanName(t *testing.T) {
	_, err := ParseTag("beanName=")
	require.ErrorIs(t, err, ErrParseTag)
}

func TestParseTag_UnknownTag(t *testing.T) {
	_, err := ParseTag("unknown")
	require.ErrorIs(t, err, ErrParseTag)
}

func TestParseTag_ExtraValues(t *testing.T) {
	_, err := ParseTag("beanName=beanA=beanB")
	require.ErrorIs(t, err, ErrParseTag)
}
