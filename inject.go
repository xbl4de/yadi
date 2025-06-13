package yadi

import (
	"github.com/pkg/errors"
	"reflect"
)

func InjectLazyBean[T interface{}]() LazyBean[T] {
	return func() T {
		return RequireBean[T]()
	}
}

func Inject(value Bean) error {
	reflectValue := reflect.ValueOf(value)
	reflectType := reflect.TypeOf(value)

	if isTypeDoesNotSupportInjection(reflectType) {
		return errors.Errorf("Inject to %s is not supported", reflectType.String())
	}

	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}

	fieldsCount := reflectType.NumField()
	for fieldNum := 0; fieldNum < fieldsCount; fieldNum++ {
		field := reflectType.Field(fieldNum)
		yadiTag, err := ParseTag(field.Tag.Get(TagName))
		if err != nil {
			return err
		}
		if yadiTag.Ignore {
			continue
		}
		fieldValue := reflectValue.Field(fieldNum)
		fieldType := fieldValue.Type()
		if fieldType.Kind() == reflect.Func {
			continue
		}
		toInject, err := getValueToInject(fieldType, yadiTag)
		if err != nil {
			return err
		}
		fieldValue.Set(reflect.ValueOf(toInject))
	}
	return nil
}

func getValueToInject(fieldType reflect.Type, yadiTag *Tag) (interface{}, error) {
	if isTypeBean(fieldType) {
		bean, err := getBean(fieldType)
		if err != nil {
			return nil, err
		}
		return bean, nil
	} else {
		path := yadiTag.ValuePath
		genericValue, err := getGenericValue(path)
		if err != nil {
			return nil, err
		}
		return genericValue, nil
	}
}

func buildLazyBean(beanType reflect.Type) func() interface{} {
	return func() interface{} {
		b, err := getBeanFromContext(beanType)
		if err != nil {
			panic(err)
		}
		return b
	}
}
