package differ

import (
	"fmt"
	"reflect"
	"strings"
)

// Diff returns a string representation of differences between two values
func Diff(expected, got interface{}) string {
	if expected == nil || got == nil {
		return fmt.Sprintf("\nExpected: %v\nGot:      %v", expected, got)
	}

	expectedType := reflect.TypeOf(expected)
	gotType := reflect.TypeOf(got)

	if expectedType != gotType {
		return fmt.Sprintf("\nType mismatch:\nExpected type: %v\nGot type:      %v", expectedType, gotType)
	}

	switch expectedType.Kind() {
	case reflect.Struct:
		return diffStruct(expected, got)
	case reflect.Map:
		return diffMap(expected, got)
	case reflect.Slice, reflect.Array:
		return diffSlice(expected, got)
	default:
		return fmt.Sprintf("\nExpected: %v\nGot:      %v", expected, got)
	}
}

func diffStruct(expected, got interface{}) string {
	var diff strings.Builder
	diff.WriteString("\n")

	expectedVal := reflect.ValueOf(expected)
	gotVal := reflect.ValueOf(got)
	typ := expectedVal.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		expectedField := expectedVal.Field(i)
		gotField := gotVal.Field(i)

		if !reflect.DeepEqual(expectedField.Interface(), gotField.Interface()) {
			diff.WriteString(fmt.Sprintf("Field %s:\n", field.Name))
			diff.WriteString(fmt.Sprintf("  Expected: %v\n", expectedField.Interface()))
			diff.WriteString(fmt.Sprintf("  Got:      %v\n", gotField.Interface()))
		}
	}

	return diff.String()
}

func diffMap(expected, got interface{}) string {
	var diff strings.Builder
	diff.WriteString("\n")

	expectedVal := reflect.ValueOf(expected)
	gotVal := reflect.ValueOf(got)

	// Check for missing keys in got
	for _, key := range expectedVal.MapKeys() {
		expectedValue := expectedVal.MapIndex(key)
		gotValue := gotVal.MapIndex(key)

		if !gotValue.IsValid() {
			diff.WriteString(fmt.Sprintf("Missing key %v in got\n", key))
			diff.WriteString(fmt.Sprintf("  Expected: %v\n", expectedValue.Interface()))
			continue
		}

		if !reflect.DeepEqual(expectedValue.Interface(), gotValue.Interface()) {
			diff.WriteString(fmt.Sprintf("Key %v:\n", key))
			diff.WriteString(fmt.Sprintf("  Expected: %v\n", expectedValue.Interface()))
			diff.WriteString(fmt.Sprintf("  Got:      %v\n", gotValue.Interface()))
		}
	}

	// Check for extra keys in got
	for _, key := range gotVal.MapKeys() {
		if !expectedVal.MapIndex(key).IsValid() {
			diff.WriteString(fmt.Sprintf("Extra key %v in got\n", key))
			diff.WriteString(fmt.Sprintf("  Got: %v\n", gotVal.MapIndex(key).Interface()))
		}
	}

	return diff.String()
}

func diffSlice(expected, got interface{}) string {
	var diff strings.Builder
	diff.WriteString("\n")

	expectedVal := reflect.ValueOf(expected)
	gotVal := reflect.ValueOf(got)

	maxLen := expectedVal.Len()
	if gotVal.Len() > maxLen {
		maxLen = gotVal.Len()
	}

	for i := 0; i < maxLen; i++ {
		if i >= expectedVal.Len() {
			diff.WriteString(fmt.Sprintf("Extra element at index %d:\n", i))
			diff.WriteString(fmt.Sprintf("  Got: %v\n", gotVal.Index(i).Interface()))
			continue
		}
		if i >= gotVal.Len() {
			diff.WriteString(fmt.Sprintf("Missing element at index %d:\n", i))
			diff.WriteString(fmt.Sprintf("  Expected: %v\n", expectedVal.Index(i).Interface()))
			continue
		}

		if !reflect.DeepEqual(expectedVal.Index(i).Interface(), gotVal.Index(i).Interface()) {
			diff.WriteString(fmt.Sprintf("Element at index %d:\n", i))
			diff.WriteString(fmt.Sprintf("  Expected: %v\n", expectedVal.Index(i).Interface()))
			diff.WriteString(fmt.Sprintf("  Got:      %v\n", gotVal.Index(i).Interface()))
		}
	}

	return diff.String()
} 