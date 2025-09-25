package main

import (
	"fmt"
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

func main() {
	val := lakery.NewValidator()
	val.RegisterTag(noSpaces, noSpacesValidator)

	s := Struct{
		Name: "killer-= feature",
	}
	if err := val.Validate(s); err != nil {
		fmt.Println(err)
	}
}
