package yadi

import (
	"github.com/pkg/errors"
	"reflect"
)

func getBeanOrDefault(beanType reflect.Type, defaultValue Bean) (Bean, error) {
	err := validateTypeIsBean(beanType)
	if err != nil {
		return nil, err
	}
	bean, err := getBeanFromContext(beanType)
	if err != nil {
		if errNoInjectableProvided(err) && defaultValue != nil {
			return defaultValue, nil
		}
		return nil, err
	}
	return bean, nil
}

func getBean(beanType reflect.Type) (Bean, error) {
	return getBeanOrDefault(beanType, nil)
}

func validateTypeIsBean(beanType reflect.Type) error {
	if beanType.Kind() == reflect.Ptr {
		if beanType.Elem().Kind() == reflect.Struct {
			return nil
		} else {
			return errors.Wrapf(ErrNonBeanType, "pointer to %s", beanType.String())
		}
	}
	if beanType.Kind() == reflect.Struct || beanType.Kind() == reflect.Interface {
		return nil
	} else {
		return errors.Wrapf(ErrNonBeanType, "%s", beanType.String())
	}
}

func isTypeBean(beanType reflect.Type) bool {
	return validateTypeIsBean(beanType) == nil
}

func isTypeSupportInjection(beanType reflect.Type) bool {
	return beanType.Kind() == reflect.Ptr && beanType.Elem().Kind() == reflect.Struct
}

func isTypeDoesNotSupportInjection(beanType reflect.Type) bool {
	return !isTypeSupportInjection(beanType)
}

func getBeanFromContext(beanType reflect.Type) (Bean, error) {
	val, err := globalCtx.Get(beanType)
	if err != nil {
		return nil, err
	}
	return val, err
}

func getGenericValueOrDefault(path string, defaultValue interface{}) (interface{}, error) {
	value, err := globalCtx.GetGenericValue(path)
	if err != nil {
		if errors.Is(err, ErrNoValueFound) && defaultValue != nil {
			return defaultValue, nil
		}
		return nil, err
	}
	return value, nil
}

func getGenericValue(path string) (interface{}, error) {
	return getGenericValueOrDefault(path, nil)
}

func convertToBean(reflectVal reflect.Value) Bean {
	return reflectVal.Interface()
}

func castToErr(errVal reflect.Value) error {
	if errVal.IsNil() {
		return nil
	}
	return errVal.Elem().Interface().(error)
}

func buildValueBox[T interface{}](reflectVal reflect.Value) *ValueBox[T] {
	return &ValueBox[T]{convertToBean(reflectVal).(T)}
}

func isTypesNotCompatible(declaredType reflect.Type, returnType reflect.Type) bool {
	if declaredType != returnType {
		if declaredType.Kind() == reflect.Interface && returnType.Kind() == reflect.Ptr {
			return !returnType.AssignableTo(declaredType)
		} else {
			return true
		}
	} else {
		return false
	}
}
