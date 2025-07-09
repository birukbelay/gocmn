package util

import "reflect"

func unwrapReflective(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface, reflect.Ptr:
		val2 := value.Elem()
		return unwrapReflective(val2)
	default:
		return value

	}
}

func UnwrapAny(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	j := reflect.ValueOf(value)
	out := unwrapReflective(j)
	if out.CanAddr() {
		return out.Addr().Interface()
	}
	as := reflect.New(out.Type())
	as.Elem().Set(out)
	return as.Interface()
}
