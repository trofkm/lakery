package main

import (
	"fmt"

	"github.com/trofkm/lakery"
)

type Struct struct {
	Name string `lakery:"min=10,max=15"`
}

func main() {
	val := lakery.NewValidator()

	s := Struct{
		Name: "killer-= feature",
	}
	if err := val.Validate(s); err != nil {
		fmt.Println(err)
	}
}
