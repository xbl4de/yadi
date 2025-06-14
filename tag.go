package yadi

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Tag struct {
	Ignore    bool
	BeanName  string
	ValuePath string
}

const TagName = "yadi"

const (
	IgnoreValue  = "ignore"
	BeanNameTag  = "beanName"
	ValuePathTag = "path"
)

type tagModifier func(*Tag, string) error

var emptyTag = Tag{}

var modifiers = map[string]tagModifier{
	IgnoreValue:  applyIgnoreTag,
	BeanNameTag:  applyBeanNameTag,
	ValuePathTag: applyPathTag,
}

func applyIgnoreTag(tag *Tag, _ string) error {
	tag.Ignore = true
	return nil
}

func applyBeanNameTag(tag *Tag, value string) error {
	if value == "" {
		return errors.Errorf("Expected non-empty beanName, but got %s", value)
	}
	tag.BeanName = value
	return nil
}

func applyPathTag(tag *Tag, value string) error {
	if value == "" {
		return errors.Errorf(`"%s" is empty`, BeanNameTag)
	}
	tag.ValuePath = value
	return nil
}

func ParseTag(tag string) (*Tag, error) {
	if strings.TrimSpace(tag) == "" {
		return &emptyTag, nil
	}
	var err error
	parts := strings.Split(tag, ";")
	newTag := &Tag{}
	for _, part := range parts {
		err = applyPart(part, newTag)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrParseTag, err)
		}
	}
	return newTag, nil
}

func applyPart(part string, tag *Tag) error {
	partArray := strings.Split(part, "=")
	if len(partArray) > 2 {
		return fmt.Errorf("expected key=value, got %s", part)
	}
	key := partArray[0]
	value := ""
	if len(partArray) == 2 {
		value = partArray[1]
	}
	modifier, ok := modifiers[key]
	if !ok {
		return errors.Errorf("unknown tag '%s'", key)
	}
	return modifier(tag, value)
}
