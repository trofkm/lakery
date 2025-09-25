package main

import (
	"fmt"
	"strconv"

	"github.com/trofkm/lakery"
)

type Struct struct {
	Name string `lakery:"min=10,max=15"`
}

const (
	min = "min"
	max = "max"
)

func minValidator(lk *lakery.Value) error {
	val := lk.String()
	minLen, err := strconv.Atoi(lk.Param())
	if err != nil {
		return err
	}
	if len(val) < minLen {
		return fmt.Errorf("min len should be at least %d", minLen)
	}
	return nil
}
func maxValidator(lk *lakery.Value) error {
	val := lk.String()
	maxLen, err := strconv.Atoi(lk.Param())
	if err != nil {
		return err
	}
	if len(val) > maxLen {
		return fmt.Errorf("max len should be less equal then %d", maxLen)
	}
	return nil
}

func main() {
	val := lakery.NewValidator()
	val.RegisterTag(min, minValidator)
	val.RegisterTag(max, maxValidator)

	s := Struct{
		Name: "killer-= feature",
	}
	if err := val.Validate(s); err != nil {
		fmt.Println(err)
	}
}
