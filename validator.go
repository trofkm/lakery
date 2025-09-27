package lakery

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	mainTag = "lakery"
)

type TagValidationFunc = func(*Value) error

type Validator struct {
	validators map[string]TagValidationFunc
}

func NewValidator() *Validator {
	v := &Validator{
		validators: make(map[string]TagValidationFunc),
	}
	// register built-in validators
	v.registerBuiltins()
	return v
}

func (v *Validator) RegisterTag(tag string, fn TagValidationFunc) {
	// don't check for existens - its totally fine to override some validator
	v.validators[tag] = fn
}

func (v *Validator) ListValidators() []string {
	// cache? not necessary since it is probably not very often to call
	vals := make([]string, 0, len(v.validators))
	for k := range v.validators {
		vals = append(vals, k)
	}
	return vals
}

func (v *Validator) Validate(s any) error {
	// todo: parse internal structure here and search for data
	if v == nil {
		return errors.New("cannot validate nil")
	}

	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errors.New("can only validate structs")
	}
	return v.validateStruct(rv)
}

func (v *Validator) validateStruct(rv reflect.Value) error {
	typ := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := typ.Field(i)
		if err := v.proceedTags(field, fieldType); err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) proceedTags(fieldValue reflect.Value, fieldType reflect.StructField) error {
	// "lakery:..." tag
	rootTag := fieldType.Tag.Get(mainTag)
	if rootTag == "" {
		return nil
	}

	tags, err := splitTopLevelByComma(rootTag)
	if err != nil {
		return CurrentErrorFormatFunc(fieldType, fieldValue, err)
	}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		var val *Value = &Value{val: fieldValue, name: fieldType.Name}
		// if we have param - put it into Value field
		splitted := strings.SplitN(tag, "=", 2)
		tagKey := strings.TrimSpace(splitted[0])

		withVal := len(splitted) == 2

		if withVal {
			val.param = strings.TrimSpace(splitted[1])
		}

		// special handling for each={...}
		if tagKey == eachTag {
			// only applicable to slices/arrays
			kind := fieldValue.Kind()
			if kind != reflect.Slice && kind != reflect.Array {
				return CurrentErrorFormatFunc(fieldType, fieldValue, fmt.Errorf("each can be used only with slice or array"))
			}
			inner := val.Param()
			inner = strings.TrimSpace(inner)
			if strings.HasPrefix(inner, "{") && strings.HasSuffix(inner, "}") {
				inner = strings.TrimSpace(inner[1 : len(inner)-1])
			}
			innerTags, err := splitTopLevelByComma(inner)
			if err != nil {
				return CurrentErrorFormatFunc(fieldType, fieldValue, err)
			}
			for i := 0; i < fieldValue.Len(); i++ {
				elem := fieldValue.Index(i)
				for _, it := range innerTags {
					it = strings.TrimSpace(it)
					if it == "" {
						continue
					}
					kv := strings.SplitN(it, "=", 2)
					innerKey := strings.TrimSpace(kv[0])
					eVal := &Value{val: elem, name: fieldType.Name}
					if len(kv) == 2 {
						eVal.param = strings.TrimSpace(kv[1])
					}
					if validator, ok := v.validators[innerKey]; ok {
						if err := validator(eVal); err != nil {
							// report error for the specific element value
							return CurrentErrorFormatFunc(fieldType, elem, err)
						}
					}
				}
			}
			continue
		}

		if validator, ok := v.validators[tagKey]; ok {
			if err := validator(val); err != nil {
				return CurrentErrorFormatFunc(fieldType, fieldValue, err)
			}
		}
	}
	return nil
}

// splitTopLevelByComma splits a string by commas, ignoring commas inside curly braces.
func splitTopLevelByComma(s string) ([]string, error) {
	var parts []string
	depth := 0
	last := 0
	for i, r := range s {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[last:i])
				last = i + 1
			}
		}
	}
	parts = append(parts, s[last:])
	if depth < 0 {
		return nil, fmt.Errorf("unopened braces in %q", s)
	} else if depth > 0 {
		return nil, fmt.Errorf("unclosed braces in %q", s)
	}
	return parts, nil
}
