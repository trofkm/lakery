package lakery

import (
	"fmt"
	"reflect"
)

// ErrorFunc used to create validation error
type ErrorFormatFunc = func(fieldType reflect.StructField, fieldValue reflect.Value, err error) error

func defaultErrorFormat(fieldType reflect.StructField, fieldValue reflect.Value, err error) error {
	return fmt.Errorf("field %q validation error: %w (received: '%v')", fieldType.Name, err, fieldValue)
}

var CurrentErrorFormatFunc ErrorFormatFunc = defaultErrorFormat
