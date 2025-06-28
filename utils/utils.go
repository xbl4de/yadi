package utils

import (
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/types"
	"reflect"
)

func ValidateTypeIsBean(beanType reflect.Type) error {
	if beanType.Kind() == reflect.Ptr {
		if beanType.Elem().Kind() == reflect.Struct {
			return nil
		} else {
			return errors.Wrapf(types.ErrNonBeanType, "pointer to %s", beanType.String())
		}
	}
	if beanType.Kind() == reflect.Struct || beanType.Kind() == reflect.Interface {
		return nil
	} else {
		return errors.Wrapf(types.ErrNonBeanType, "%s", beanType.String())
	}
}

func IsTypeBean(beanType reflect.Type) bool {
	return ValidateTypeIsBean(beanType) == nil
}

func IsTypeSupportInjection(beanType reflect.Type) bool {
	return beanType.Kind() == reflect.Ptr && beanType.Elem().Kind() == reflect.Struct
}

func IsTypeDoesNotSupportInjection(beanType reflect.Type) bool {
	return !IsTypeSupportInjection(beanType)
}

func ConvertToBean(reflectVal reflect.Value) types.Bean {
	return reflectVal.Interface()
}

func CastToErr(errVal reflect.Value) error {
	if errVal.IsNil() {
		return nil
	}
	return errVal.Elem().Interface().(error)
}

func IsTypesNotCompatible(declaredType reflect.Type, returnType reflect.Type) bool {
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
