// Package config implements utilities for application configuration.
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

// ValidateConstraints checks the struct tags on the provided struct, returning an error if any constraint fails.
// An error is also returned if: the provided value is nil, or not a struct or a ptr-to-struct; if the struct contains
// unsupported data types; or if a malformed constraint is found.
//
// Supported struct tags:
// 		config:"required" -- the value must not be the zero-value for the type.
//
// Example:
//		type Config struct {
//			MyIntValue int `config:"required"` // must not equal 0
//			MyStringValue string `config:"required"` // must not equal ""
//		}
func ValidateConstraints(config interface{}) error {
	return validateStruct("", config)
}

func isZeroValue(valueType reflect.Type, value reflect.Value) bool {

	if !value.IsValid() { // zero Value type -- nil in caller passed to reflect.TypeOf/reflect.ValueOf that was then passed to us
		return true
	}

	typeKind := valueType.Kind()

	switch typeKind {
	// re: including reflect.Ptr in this list: we want to check if the ptr is nil, NOT whether the thing it points to,
	// if anything, is a zero-value
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()

	case reflect.Array:
		for i := 0; i < value.Len(); i++ {
			elementValue := value.Index(i)
			if !isZeroValue(elementValue.Type(), elementValue) {
				return false
			}
		}
		return true

	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)

			// skip unexported fields -- go panics if we try to return a value obtained from an unexported field (via
			// Interface(); if we pass an unexported field to isZeroValue and it is a default type that hits the
			// Interface() call at the end of this method, then we experience the panic).
			if fieldValue.CanInterface() {
				if !isZeroValue(fieldValue.Type(), fieldValue) {
					return false
				}
			}
		}
		return true

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

			switch fieldType.Type.Kind() {
			// we'll skip supporting a number of types that don't make sense / don't have a usecase at the moment
			case reflect.Bool, reflect.Uintptr, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
				reflect.Map, reflect.Slice, reflect.UnsafePointer:
				return fmt.Errorf("%s: required field is an unsupported type kind of %v (%v)", fieldNameWithPrefix(prefix, fieldType.Name), fieldType.Type.Kind(), fieldType.Type)
			}

			if isZeroValue(fieldType.Type, fieldValue) {
				return fmt.Errorf("%s: required field of type %v (kind %v) has a zero-value", fieldNameWithPrefix(prefix, fieldType.Name), fieldType.Type, fieldType.Type.Kind())
			}
		} else {
			return fmt.Errorf("%s: unknown %s tag value \"%s\"", fieldNameWithPrefix(prefix, fieldType.Name), ConfigTag, tagValue)
		}
	}

	return nil
}

func validateStruct(prefix string, any interface{}) error {

	if any == nil {
		return fmt.Errorf("%s: zero-value", prefix)
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

		default:
			if err := examineFieldTags(prefix, fieldType, fieldValue); err != nil {
				return err
			}

			// Recurse into ptr-to-structs (embedded value structs will have been detected above)
			if typeKind == reflect.Ptr && !fieldValue.IsNil() && fieldType.Type.Elem().Kind() == reflect.Struct {
				if err := validateStruct(fieldNameWithPrefix(prefix, fieldType.Name), fieldValue.Interface()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
