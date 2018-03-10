package config

import (
	"fmt"
	"reflect"
)

// ConfigTag is the name of the config struct tag.
// e.g. `config:"required"`
const ConfigTag = "config"

// RequiredTagValue is the name of 'required' option of the config struct tag.
// e.g. `config:"required"`
const RequiredTagValue = "required"

// Configer is implemented by configurations.
type Configer interface {

	// IsValid returns an error if the current configuration is invalid for any reason.
	// This function is typically called after configuration settings have been loaded into a Configer.
	IsValid() error
}

// ValidateConstraints checks the struct tags on the provided Configer, returning an error if any constraint fails.
// This function should be used as part of your Configer's IsValid method.
func ValidateConstraints(config Configer) error {
	return validateStruct("config", config)
}

// (note that this function is non-exhaustive, e.g. arrays and structs aren't supported; see validateStruct
func isZeroValue(valueType reflect.Type, value reflect.Value) bool {

	typeKind := valueType.Kind()

	switch typeKind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()

	// These types need special handling, and we've already covered them below in validateStruct (i.e. they shouldn't be passed to this function)
	case reflect.Array, reflect.Struct:
		panic(fmt.Sprintf("shouldn't have received kind of %v: unsupported type", typeKind))
	}

	return value.Interface() == reflect.Zero(valueType).Interface()
}

func fieldNameWithPrefix(prefix string, fieldName string) string {
	if prefix != "" {
		return prefix + "." + fieldName
	}

	return fieldName
}

func examineFieldTags(prefix string, fieldType reflect.StructField, fieldValue reflect.Value) error {

	if tagValue, ok := fieldType.Tag.Lookup(ConfigTag); ok {
		if tagValue == RequiredTagValue {
			if isZeroValue(fieldType.Type, fieldValue) {
				return fmt.Errorf("%s: required field of type %s has a zero-value", fieldNameWithPrefix(prefix, fieldType.Name), fieldType.Type.Name())
			}
		} else {
			return fmt.Errorf("%s: unknown %s tag value \"%s\"", fieldNameWithPrefix(prefix, fieldType.Name), ConfigTag, tagValue)
		}
	}

	return nil
}

func validateStruct(prefix string, any interface{}) error {

	if any == nil {
		return fmt.Errorf("%s: Zero-value", prefix)
	}

	anyType := reflect.TypeOf(any)
	anyValue := reflect.ValueOf(any)

	if anyType.Kind() == reflect.Ptr {
		anyType = anyType.Elem()
		anyValue = anyValue.Elem()
	}

	if anyType.Kind() != reflect.Struct {
		return fmt.Errorf("%s: expected a struct (or ptr-to-struct) but received a kind of %v", prefix, anyType.Kind())
	}

	for i := 0; i < anyValue.NumField(); i++ {
		fieldType := anyType.Field(i)
		fieldValue := anyValue.Field(i)
		typeKind := fieldType.Type.Kind()

		switch typeKind {
		case reflect.Struct:
			if err := validateStruct(fieldNameWithPrefix(prefix, fieldType.Name), fieldValue.Interface()); err != nil {
				return err
			}

		// we'll skip supporting a number of types that don't make sense / don't have a usecase at the moment
		case reflect.Uintptr, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
			return fmt.Errorf("%s: unsupported type of kind of %v", fieldNameWithPrefix(prefix, fieldType.Name), typeKind)

		default:
			if err := examineFieldTags(prefix, fieldType, fieldValue); err != nil {
				return err
			}
		}
	}

	return nil
}
