package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/dtos"
)

func GetFieldValue(obj interface{}, fieldName string) (fieldValue string, error error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // Dereference if it's a pointer
	}

	// Traverse the fields to find the value
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Check if the field is embedded (anonymous)
		if fieldType.Anonymous {
			embeddedValue, err := GetFieldValue(field.Interface(), fieldName)
			if err == nil {
				return embeddedValue, nil
			}
		} else if fieldType.Name == fieldName {

			return findFieldValue(field)
			//return fmt.Sprintf("%v", field.Interface()), nil
		}
	}

	return "", fmt.Errorf("field %s not found", fieldName)
}

func findFieldValue(fieldValue reflect.Value) (string, error) {
	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}
	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", fieldValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", fieldValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", fieldValue.Float()), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", fieldValue.Bool()), nil
	case reflect.Struct:
		if t, ok := fieldValue.Interface().(time.Time); ok {
			return t.Format(time.RFC3339Nano), nil
		}
		bytes, err := json.Marshal(fieldValue.Interface())
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	case reflect.Slice, reflect.Array, reflect.Map:
		bytes, err := json.Marshal(fieldValue.Interface())
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	default:
		return fmt.Sprintf("%v", fieldValue.Interface()), nil
	}
}

// ToPascalCase is a helper function to convert snake_case & camelCase to PascalCase
// this is because Golang's struct names are all PascalCase
func ToPascalCase(input string) string {
	words := strings.Split(input, "_")
	for i := range words {
		words[i] = cases.Title(language.Und).String(words[i])
	}
	return strings.Join(words, "")
}

func GetPagiValues(pagi dtos.PaginationInput) (orderedBy, sortDirection string, limit, offSet int) {
	sortBy := "created_at"
	//TODO
	if pagi.SortBy != "" {
		sortBy = pagi.SortBy
	}

	sortDir := "asc"
	if pagi.SortDir != "asc" {
		// will make the default desc
		sortDir = "desc"
	}
	pagLimit := pagi.Limit
	if pagi.Limit == 0 {
		pagLimit = 20
	}
	offsets := pagi.Page - 1
	if offsets < 1 {
		offsets = 0
	}
	//page 3 means the offset is 2, page 1 is offset 0
	return sortBy, sortDir, pagLimit, offsets
}

// GenerateQueryString creates the sort direction and sort by
func GenerateQueryString(pagi dtos.PaginationInput) (orderString, paginationString, orderBy string) {
	orderedBy := "created_at"
	if pagi.SortBy != "" {
		orderedBy = pagi.SortBy
	}
	sortDirection := "desc"
	if pagi.SortDir != "" {
		sortDirection = pagi.SortDir
	}
	var queryString string
	if sortDirection == "asc" {
		queryString = fmt.Sprintf("(%s > ?) OR (%s = ? AND id > ?)", orderedBy, orderedBy)
	} else {
		queryString = fmt.Sprintf("(%s < ?) OR (%s = ? AND id > ?)", orderedBy, orderedBy)
	}
	//eg `created_at asc,id`
	sortOrder := fmt.Sprintf("%s %s, id", orderedBy, sortDirection)
	return sortOrder, queryString, orderedBy
}

// GenerateFwdQueryString creates the sort direction and sort by
func GenerateFwdQueryString(orderedBy, sortDirection string) (orderString, paginationString string) {
	var queryString string
	if sortDirection != "asc" {
		sortDirection = "desc"
	}
	if sortDirection == "asc" {
		queryString = fmt.Sprintf("(%s > ?) OR (%s = ? AND id > ?)", orderedBy, orderedBy)
	} else {
		queryString = fmt.Sprintf("(%s < ?) OR (%s = ? AND id < ?)", orderedBy, orderedBy)
	}
	//eg `created_at asc, id `
	sortOrder := fmt.Sprintf("%s %s, id %s", orderedBy, sortDirection, sortDirection)
	return sortOrder, queryString
}

// GeneratePrevQueryString creates the sort direction and sort by
func GeneratePrevQueryString(orderedBy, sortDirection string) (orderString, paginationString string) {

	var queryString string
	sortDir := "asc"
	if sortDirection == "asc" {
		sortDir = "desc"
		queryString = fmt.Sprintf("(%s < ?) OR (%s = ? AND id < ?)", orderedBy, orderedBy)
	} else {
		sortDir = "asc"
		queryString = fmt.Sprintf("(%s > ?) OR (%s = ? AND id > ?)", orderedBy, orderedBy)
	}
	//to get the prev page the id will be descending, because previously it was ascending
	sortString := fmt.Sprintf("%s %s, id %s", orderedBy, sortDir, sortDir)
	return sortString, queryString
}
func GenerateCursor(orderedBy string, lastElement any, forward bool) (cursor string, err error) {
	//we change it to PascalCase because golang's struct names are all PascalCase
	structNameVal := ToPascalCase(orderedBy)
	value, err := GetFieldValue(lastElement, structNameVal)
	if err != nil {
		return "", err
	}
	idValue, err := GetFieldValue(lastElement, "ID")
	if err != nil {
		return "", err
	}
	//f := v.FieldByName(structNameVal)
	cmn.LogTrace("the order is", value)
	cursor = EncodeCursor(orderedBy, value, idValue, forward)
	return cursor, nil
}
func DecodeCursor(opaqueCursor string) (OrderBy, cursorsValue, cursorsId *string, forward bool, err error) {
	if opaqueCursor == "" {
		return nil, nil, nil, true, nil
	}

	decodedCursor, err := base64.StdEncoding.DecodeString(opaqueCursor)
	if err != nil {
		return nil, nil, nil, true, err
	}
	var cursorVal, cursorID, OrderByName string
	parts := strings.Split(string(decodedCursor), "|")
	if len(parts) == 4 {
		OrderByName = parts[0]
		cursorVal = parts[1]
		cursorID = parts[2]
		fwd, err := strconv.ParseBool(parts[3])
		if err != nil {
			return nil, nil, nil, true, err
		}
		return &OrderByName, &cursorVal, &cursorID, fwd, nil
	}
	return nil, nil, nil, true, fmt.Errorf("invalid cursor format")

}

func EncodeCursor(cursorName, cursorValue, id string, forward bool) string {
	cursorString := fmt.Sprintf("%s|%s|%s|%t", cursorName, cursorValue, id, forward)
	return base64.StdEncoding.EncodeToString([]byte(cursorString))

}
