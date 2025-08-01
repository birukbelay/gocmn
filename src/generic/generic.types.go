package generic

import (
	"reflect"

	"gorm.io/gorm/clause"
)

type KeyArr map[string][]string

type Opt struct {
	Debug         bool
	Preloads      []string
	NoLimit       bool
	InQueries     KeyArr
	UpdateColumns []string
	Columns       []clause.Column
	AuthKey       *string
	AuthVal       *string
}

type AssociationKey string

const (
	AssociationName = AssociationKey("name")
	AssociationId   = AssociationKey("id")
	AssociationSlug = AssociationKey("slug")
)

type AssocVar struct {
	ModelName           string
	Key                 AssociationKey
	AssociatedValues    []string
	EmptyingAssociation bool
	Debug               bool
	Preloads            []string
	AuthKey             *string
	AuthVal             *string
}

func isEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)
	// If the value is a pointer, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// Check if it's a struct now
	if v.Kind() != reflect.Struct {
		return false // Not a struct, return false
	}
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			return false
		}
	}
	return true
}

type Queryable interface {
	GetQueries() (string, []string)
}
