package handlers

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetJSONFieldName[T any](fe validator.FieldError, obj T) string {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Look up the field by its struct name as given by the FieldError.
	field, found := t.FieldByName(fe.StructField())
	if !found {
		// Fallback to using the field name from the validator.
		return fe.Field()
	}
	jsonTag := field.Tag.Get("json")
	// If the json tag is empty or ignored, fallback.
	if jsonTag == "" || jsonTag == "-" {
		return fe.Field()
	}
	// If the tag contains comma options (e.g. "question,omitempty"), return only the first part.
	return strings.Split(jsonTag, ",")[0]
}
