package lakery

import "reflect"

type Value struct {
	val reflect.Value
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
