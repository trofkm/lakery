package lakery

import (
	"fmt"
	"reflect"
)

type Value struct {
	val   reflect.Value
	name  string
	param string
}

// todo: this is very interesting question - how we can obtain the underlaying value
func (v *Value) String() string {
	return v.val.String()
}

func (v *Value) Interface() any {
	if v.val.CanInterface() {
		return v.val.Interface()
	}
	panic(v.val.Type().String() + " is not an interface type")
}

func (v *Value) Param() string {
	if v.param != "" {
		return v.param
	}
	panic(fmt.Sprintf("requested param value for %q is not set", v.name))
}
