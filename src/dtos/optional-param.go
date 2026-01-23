package dtos

import (
	"database/sql/driver"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type OptParam[T any] struct {
	Val   T
	IsSet bool
}

// Define schema to use wrapped type
func (o OptParam[T]) Schema(r huma.Registry) *huma.Schema {
	return huma.SchemaFromType(r, reflect.TypeFor[T]())
}

// Receiver tells Huma where to store the parsed value.
func (o *OptParam[T]) Receiver() reflect.Value {
	v := reflect.ValueOf(o).Elem().FieldByName("Val")
	return v
}

// React to request param being parsed to update internal state
// MUST have pointer receiver
func (o *OptParam[T]) OnParamSet(isSet bool, parsed any) {
	o.IsSet = isSet
}

// Value implements the driver.Valuer interface for InfoMap
// this MUST NOT have pointer receiver
func (o OptParam[T]) Value() (driver.Value, error) {
	// Convert InfoMap to JSON
	if !o.IsSet {
		return nil, nil
	}
	if t, ok := any(o.Val).(driver.Valuer); ok {
		return t.Value()
	}
	return o.Val, nil
}

// Receiver tells Huma where to store the parsed value.
func NewOptonal[T any](val T) OptParam[T] {
	return OptParam[T]{Val: val, IsSet: true}
}
