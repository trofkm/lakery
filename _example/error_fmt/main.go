package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/trofkm/lakery"
)

type Struct struct {
	Name string `lakery:"no_spaces"`
}

const noSpaces = "no_spaces"

func noSpacesValidator(lk *lakery.Value) error {
	val := lk.String()
	if strings.Contains(val, " ") {
		return fmt.Errorf("should not contain spaces")
	}
	return nil
}

// example of custom error formatter
func ruErrorFormat(fieldType reflect.StructField, fieldValue reflect.Value, err error) error {
	return fmt.Errorf("ошибка валидации: '%w' для поля %q  (получено: '%v')", err, fieldType.Name, fieldValue)
}

func main() {
	lakery.CurrentErrorFormatFunc = ruErrorFormat
	val := lakery.NewValidator()
	val.RegisterTag(noSpaces, noSpacesValidator)

	s := Struct{
		Name: "killer-= feature",
	}
	if err := val.Validate(s); err != nil {
		fmt.Println(err)
	}
}
