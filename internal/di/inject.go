package di

import (
	"github.com/pkg/errors"
	"github.com/xbl4de/yadi/internal/log"
	"github.com/xbl4de/yadi/internal/utils"
	"github.com/xbl4de/yadi/types"
	"reflect"
)

func InjectLazyBean[T types.Bean]() types.LazyBean[T] {
	return func() T {
		return RequireBean[T]()
	}
}

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

func Inject(valuePtr types.Bean) error {
	return injectToPtr(reflect.ValueOf(valuePtr))
}

func injectToPtr(reflectValue reflect.Value) error {
	reflectType := reflectValue.Type()

	if utils.IsTypeDoesNotSupportInjection(reflectType) {
		return errors.Wrapf(types.ErrInjectNotSupported, "Inject %s to value failed", reflectType.String())
	}

	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}

	fieldsCount := reflectType.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := reflectType.Field(i)
		fieldValue := reflectValue.Field(i)
		log.Verbose("Injecting field %s.%s", reflectType.String(), field.Name)
		err := injectToField(field, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func injectToField(field reflect.StructField, fieldValue reflect.Value) error {
	yadiTag, err := types.ParseTag(field.Tag.Get(types.TagName))
	if err != nil {
		return err
	}
	fieldType := fieldValue.Type()
	if shouldIgnoreInjection(yadiTag, fieldType) {
		return nil
	}
	toInject, err := getValueToInject(fieldType, yadiTag)
	if err != nil {
		return err
	}
	fieldValue.Set(reflect.ValueOf(toInject))
	return nil
}

func shouldIgnoreInjection(yadiTag *types.Tag, fieldType reflect.Type) bool {
	return yadiTag.Ignore || fieldType.Kind() == reflect.Func
}

func getValueToInject(fieldType reflect.Type, yadiTag *types.Tag) (interface{}, error) {
	if utils.IsTypeBean(fieldType) {
		bean, err := GetBeanFromContext(fieldType)
		if err != nil {
			return nil, err
		}
		return bean, nil
	} else {
		path := yadiTag.ValuePath
		genericValue, err := GetGenericValue(path)
		if err != nil {
			return nil, err
		}
		return genericValue, nil
	}
}
