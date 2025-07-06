package yadi

import (
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/log"
	"github.com/xbl4de/yadi/types"
	"github.com/xbl4de/yadi/utils"
	"reflect"
)

func tryToBuildNewBean(beanType reflect.Type) (interface{}, error) {
	log.Verbose("Trying to build new bean for type %s", beanType.String())
	err := utils.ValidateTypeIsBean(beanType)
	if err != nil {
		return nil, err
	}

	isTargetTypeIsPointer := beanType.Kind() == reflect.Ptr

	buildType := beanType
	if isTargetTypeIsPointer {
		buildType = buildType.Elem()
	}

	valPtr := reflect.New(buildType)
	err = injectToPtr(valPtr)
	if err != nil {
		return nil, err
	}

	log.Verbose("Built new bean for type %s", beanType.String())
	if isTargetTypeIsPointer {
		return valPtr.Interface(), nil
	} else {
		return valPtr.Elem().Interface(), nil
	}
}

func injectToPtr(beanStructValue reflect.Value) error {
	beanStructType := beanStructValue.Type()

	if utils.IsTypeDoesNotSupportInjection(beanStructType) {
		return errors.Wrapf(types.ErrInjectNotSupported, "Inject %s to value failed", beanStructType.String())
	}
	origBeanReflectValue := beanStructValue
	origBeanTypeValue := beanStructType

	if beanStructType.Kind() == reflect.Ptr {
		beanStructType = beanStructType.Elem()
		beanStructValue = beanStructValue.Elem()
	}

	fieldsCount := beanStructType.NumField()
	for i := 0; i < fieldsCount; i++ {
		err := setField(i, beanStructValue, beanStructType, origBeanTypeValue, origBeanReflectValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(
	fieldInd int,
	beanStructValue reflect.Value,
	beanStructType reflect.Type,
	origBeanTypeValue reflect.Type,
	origBeanReflectValue reflect.Value,
) error {
	field := beanStructType.Field(fieldInd)
	fieldValue := beanStructValue.Field(fieldInd)

	yadiTag, err := types.ParseTag(field.Tag.Get(types.TagName))
	if err != nil {
		return err
	}
	if shouldIgnoreInjection(yadiTag, field.Type) {
		return nil
	}
	toInject, err := getValueToInject(field.Type, yadiTag)
	if err != nil {
		return err
	}

	log.Verbose("Injecting field %s.%s", beanStructType.String(), field.Name)
	if fieldValue.CanSet() {
		err := injectToField(fieldValue, toInject)
		if err != nil {
			return err
		}
	} else {
		err := tryToUseSetter(origBeanTypeValue, origBeanReflectValue, field, toInject)
		if err != nil {
			return err
		}
	}
	return nil
}

func injectToField(field reflect.Value, value interface{}) error {
	field.Set(reflect.ValueOf(value))
	return nil
}

func tryToUseSetter(
	beanType reflect.Type,
	beanValue reflect.Value,
	reflectField reflect.StructField,
	fieldValue interface{},
) error {
	setterName := "Set" + utils.Capitalize(reflectField.Name)
	method, err := validateSetterMethod(beanType, setterName, reflectField)
	if err != nil {
		return err
	}
	method.Func.Call([]reflect.Value{beanValue, reflect.ValueOf(fieldValue)})
	return nil
}

func validateSetterMethod(beanType reflect.Type, setterName string, reflectField reflect.StructField) (reflect.Method, error) {
	method, ok := beanType.MethodByName(setterName)
	var zeroMethod reflect.Method
	if !ok {
		return zeroMethod, errors.Wrapf(types.ErrInjectNotSupported, "field %s has no setter", reflectField.Name)
	}
	if method.Type.NumIn() != 2 {
		return zeroMethod, errors.Wrapf(types.ErrInjectNotSupported, "field %s has wrong setter signature", reflectField.Name)
	}
	if method.Type.NumOut() > 0 {
		log.Verbose("Waring: %s return values will be ignored", setterName)
	}
	if !reflectField.Type.AssignableTo(method.Type.In(1)) {
		return zeroMethod, errors.Wrapf(types.ErrInjectNotSupported, "field %s has not assignable to %s", reflectField.Name, beanType.String())
	}
	return method, nil
}

func shouldIgnoreInjection(yadiTag *types.Tag, fieldType reflect.Type) bool {
	return yadiTag.Ignore || fieldType.Kind() == reflect.Func
}

func getValueToInject(fieldType reflect.Type, yadiTag *types.Tag) (interface{}, error) {
	if utils.IsTypeBean(fieldType) {
		bean, err := getBeanFromContext(fieldType)
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
