package main

import (
	"fmt"
	"strings"

	"github.com/trofkm/lakery"
)

type Struct struct {
	Name string `lakery:"min=0,max=255"`
}

const min = "min"

func minValidator(lk *lakery.Value) error {
	val := lk.String()
	if len(val)< lk.Param.(int)
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
