package base

import (
	"reflect"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Base struct {
	ID        string     `msgpack:"-" gorm:"primarykey" json:"id,omitempty"`
	CreatedAt *time.Time `msgpack:"-" json:"created_at,omitempty" gorm:"index"`
	UpdatedAt *time.Time `msgpack:"-" json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

func (m *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = ulid.Make().String()
	}
	return nil
}
func (u *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = Ptr(time.Now())
	return nil
}

func Ptr[T any](t T) *T {
	return &t
}

func SumField(n any) int64 {
	var sum int64
	v := reflect.ValueOf(n)

	// Iterate over all fields of the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Int64 {
			sum += field.Int()
		}
	}
	return sum
}
