package lakery

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	// min value for numbers, size for arrays and strings
	minTag = "min"
	// max value for numbers, size for arrays and strings
	maxTag = "max"
	// special tag for specifying validation rules for values in arrays
	eachTag = "each"
	// special tag for diving into struct type inside structure
	diveTag = "dive"
	// special tag for required fields
	requiredTag = "required"
)

// registerBuiltins registers built-in validators into the provided validator instance.
// Built-ins: min, max, required. Special tags: each, dive are handled in tag processing flow.
func (v *Validator) registerBuiltins() {
	v.RegisterTag(minTag, builtinMin)
	v.RegisterTag(maxTag, builtinMax)
	v.RegisterTag(requiredTag, builtinRequired)
}

// builtinMin validates that a value is not less than the provided minimum.
// - strings, arrays, slices: len(value) >= min
// - integers (signed/unsigned) and floats: value >= min
func builtinMin(val *Value) error {
	minStr := val.Param()
	minInt, err := strconv.Atoi(minStr)
	if err != nil {
		return fmt.Errorf("min expects integer param: %w", err)
	}
	min := minInt

	rv := val.val
	k := rv.Kind()
	if k == reflect.Pointer {
		if rv.IsNil() {
			// nil pointer fails len-based checks; consider nil < min unless min <= 0
			if min > 0 {
				return fmt.Errorf("should have length at least %d", min)
			}
			return nil
		}
		rv = rv.Elem()
		k = rv.Kind()
	}

	switch k {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if rv.Len() < min {
			return fmt.Errorf("should have length at least %d", min)
		}
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Int() < int64(min) {
			return fmt.Errorf("should be >= %d", min)
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if rv.Uint() < uint64(min) {
			return fmt.Errorf("should be >= %d", min)
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if rv.Float() < float64(min) {
			return fmt.Errorf("should be >= %d", min)
		}
		return nil
	default:
		return fmt.Errorf("min is not applicable to type %s", rv.Type())
	}
}

// builtinMax validates that a value is not greater than the provided maximum.
// - strings, arrays, slices: len(value) <= max
// - integers (signed/unsigned) and floats: value <= max
func builtinMax(val *Value) error {
	maxStr := val.Param()
	maxInt, err := strconv.Atoi(maxStr)
	if err != nil {
		return fmt.Errorf("max expects integer param: %w", err)
	}
	max := maxInt

	rv := val.val
	k := rv.Kind()
	if k == reflect.Pointer {
		if rv.IsNil() {
			// nil pointer has length 0; only passes if max >= 0
			return nil
		}
		rv = rv.Elem()
		k = rv.Kind()
	}

	switch k {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if rv.Len() > max {
			return fmt.Errorf("should have length at most %d", max)
		}
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Int() > int64(max) {
			return fmt.Errorf("should be <= %d", max)
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if rv.Uint() > uint64(max) {
			return fmt.Errorf("should be <= %d", max)
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if rv.Float() > float64(max) {
			return fmt.Errorf("should be <= %d", max)
		}
		return nil
	default:
		return fmt.Errorf("max is not applicable to type %s", rv.Type())
	}
}

// builtinRequired validates that a value is not the zero value (non-empty string, non-zero number,
// non-nil pointer/slice/map/function/interface, and structs with any non-zero field).
func builtinRequired(val *Value) error {
	rv := val.val
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return fmt.Errorf("is required")
		}
		rv = rv.Elem()
	}
	if rv.IsZero() {
		return fmt.Errorf("is required")
	}
	return nil
}
