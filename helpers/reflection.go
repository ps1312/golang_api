package helpers

import "reflect"

// GetFieldFromStruct helper to reduce code duplication
func GetFieldFromStruct(subject interface{}, field string) string {
	iter := reflect.ValueOf(subject)
	str := iter.FieldByName(field).String()
	return str
}
