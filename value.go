package lakery

import "reflect"

type Value struct {
	val   reflect.Value
	param interface{}
}

// todo: this is very interesting question - how we can obtain the underlaying value
func (v *Value) String() string {
	return v.val.String()
}

func (v *Value) Interface() interface{} {
	if v.val.CanInterface() {
		return v.val.Interface()
	}
	panic(v.val.Type().String() + " is not an interface type")
}

func (v *Value) Param() interface{} {
	if v.param != nil {
		return v.param
	}
	panic("requested param value for " + v.val.Type().String() + " is nil")
}
