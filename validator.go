package lakery

import (
	"errors"
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
	return &Validator{
		validators: make(map[string]TagValidationFunc),
	}
}

func (v *Validator) RegisterTag(tag string, fn TagValidationFunc) {
	// don't check for existens - its totally fine to override some validator
	v.validators[tag] = fn
}

func (v *Validator) ListValidators() []string {
	// cache? not necessary since it is probably not very often to call
	vals := make([]string, 0, len(v.validators))
	for k, _ := range v.validators {
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

	tags := strings.SplitSeq(rootTag, ",")
	for tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		var val *Value = &Value{val: fieldValue, name: fieldType.Name}
		// if we have param - put it into Value field
		splitted := strings.Split(tag, "=")
		tag = splitted[0]

		withVal := len(splitted) == 2

		if withVal {
			val.param = splitted[1]
		}

		if validator, ok := v.validators[tag]; ok {
			if err := validator(val); err != nil {
				return CurrentErrorFormatFunc(fieldType, fieldValue, err)
			}
		}
	}
	return nil
}
