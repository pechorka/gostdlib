package require

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/testing/differ"
)

// Equal is a helper for comparing 2 variables in non deep-equal way
func Equal[T comparable](t testing.TB, expect, got T) {
	t.Helper()
	equal(t, expect, got, "Values are not equal")
}

// Equalx is Equal with custom error message
func Equalx[T comparable](t testing.TB, expect, got T, format string, args ...any) {
	t.Helper()
	equal(t, expect, got, fmt.Sprintf(format, args...))
}

func equal[T comparable](t testing.TB, expect, got T, message string) {
	if expect != got {
		t.Logf("%s\n%s", message, differ.Diff(expect, got))
		t.FailNow()
	}
}

// EqualValues compares two values that can't be compared by == operator
func EqualValues(t testing.TB, expect, got any) {
	t.Helper()
	equalValues(t, expect, got, "Values are not equal")
}

// EqualValuesx is EqualValues with custom error message
func EqualValuesx(t testing.TB, expect, got any, format string, args ...any) {
	t.Helper()
	equalValues(t, expect, got, fmt.Sprintf(format, args...))
}

func equalValues(t testing.TB, expect, got any, message string) {
	if !reflect.DeepEqual(expect, got) {
		t.Logf("%s\n%s", message, differ.Diff(expect, got))
		t.FailNow()
	}
}

// NotEqual is a helper for ensuring that 2 variables are not equal
func NotEqual[T comparable](t testing.TB, expect, got T) {
	t.Helper()
	notEqual(t, expect, got, "Values should not be equal")
}

// NotEqualx is NotEqual with custom error message
func NotEqualx[T comparable](t testing.TB, expect, got T, format string, args ...any) {
	t.Helper()
	notEqual(t, expect, got, fmt.Sprintf(format, args...))
}

func notEqual[T comparable](t testing.TB, expect, got T, message string) {
	if expect == got {
		t.Fatalf("%s: both values are %v", message, expect)
	}
}

// NotEqualValues is a helper for ensuring that 2 values are not equal using deep comparison
func NotEqualValues(t testing.TB, expect, got any) {
	t.Helper()
	notEqualValues(t, expect, got, "Values should not be equal")
}

// NotEqualValuesx is NotEqualValues with custom error message
func NotEqualValuesx(t testing.TB, expect, got any, format string, args ...any) {
	t.Helper()
	notEqualValues(t, expect, got, fmt.Sprintf(format, args...))
}

func notEqualValues(t testing.TB, expect, got any, message string) {
	if reflect.DeepEqual(expect, got) {
		t.Fatalf("%s: both values are %v", message, expect)
	}
}

// Contains is a helper for ensuring that a string contains a substring
func Contains(t testing.TB, s, substr string) {
	t.Helper()
	contains(t, s, substr, fmt.Sprintf("String %q does not contain %q", s, substr))
}

// Containsx is Contains with custom error message
func Containsx(t testing.TB, s, substr string, format string, args ...any) {
	t.Helper()
	contains(t, s, substr, fmt.Sprintf(format, args...))
}

func contains(t testing.TB, s, substr string, message string) {
	if !strings.Contains(s, substr) {
		t.Fatal(message)
	}
}

// NotContains is a helper for ensuring that a string does not contain a substring
func NotContains(t testing.TB, s, substr string) {
	t.Helper()
	notContains(t, s, substr, fmt.Sprintf("String %q should not contain %q", s, substr))
}

// NotContainsx is NotContains with custom error message
func NotContainsx(t testing.TB, s, substr string, format string, args ...any) {
	t.Helper()
	notContains(t, s, substr, fmt.Sprintf(format, args...))
}

func notContains(t testing.TB, s, substr string, message string) {
	if strings.Contains(s, substr) {
		t.Fatal(message)
	}
}

// Empty is a helper for ensuring that a value is empty
func Empty(t testing.TB, value any) {
	t.Helper()
	empty(t, value, "Expected empty value")
}

// Emptyx is Empty with custom error message
func Emptyx(t testing.TB, value any, format string, args ...any) {
	t.Helper()
	empty(t, value, fmt.Sprintf(format, args...))
}

func empty(t testing.TB, value any, message string) {
	if !isEmpty(value) {
		t.Fatalf("%s, got: %v", message, value)
	}
}

// NotEmpty is a helper for ensuring that a value is not empty
func NotEmpty(t testing.TB, value any) {
	t.Helper()
	notEmpty(t, value, "Expected non-empty value")
}

// NotEmptyx is NotEmpty with custom error message
func NotEmptyx(t testing.TB, value any, format string, args ...any) {
	t.Helper()
	notEmpty(t, value, fmt.Sprintf(format, args...))
}

func notEmpty(t testing.TB, value any, message string) {
	if isEmpty(value) {
		t.Fatal(message)
	}
}

// isEmpty checks if value is empty (nil, zero length, or zero value)
func isEmpty(value any) bool {
	if isNil(value) {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		// Check if struct is zero value
		zero := reflect.Zero(v.Type())
		return reflect.DeepEqual(v.Interface(), zero.Interface())
	default:
		return false
	}
}

// Error is a helper for ensuring that error is not nil
func Error(t testing.TB, err error) {
	t.Helper()
	requireError(t, err, "Expected error but got nil")
}

// Errorx is Error with custom error message
func Errorx(t testing.TB, err error, format string, args ...any) {
	t.Helper()
	requireError(t, err, fmt.Sprintf(format, args...))
}

func requireError(t testing.TB, err error, message string) {
	if err == nil {
		t.Fatal(message)
	}
}

// ErrorIs is a helper for ensuring that error is target error
func ErrorIs(t testing.TB, err error, target error) {
	t.Helper()
	errorIs(t, err, target, fmt.Sprintf("Expected error %v, got %v", target, err))
}

// ErrorIsx is ErrorIs with custom error message
func ErrorIsx(t testing.TB, err error, target error, format string, args ...any) {
	t.Helper()
	errorIs(t, err, target, fmt.Sprintf(format, args...))
}

func errorIs(t testing.TB, err error, target error, message string) {
	if !errs.Is(err, target) {
		t.Fatal(message)
	}
}

// NoError is a helper for ensuring that error is nil
func NoError(t testing.TB, err error) {
	t.Helper()
	noError(t, err, fmt.Sprintf("Expected no error, got: %v", err))
}

// NoErrorx is NoError with custom error message
func NoErrorx(t testing.TB, err error, format string, args ...any) {
	t.Helper()
	noError(t, err, fmt.Sprintf(format, args...))
}

func noError(t testing.TB, err error, message string) {
	if err != nil {
		t.Fatal(message)
	}
}

// Nil is a helper for ensuring that value is nil
func Nil(t testing.TB, value any) {
	t.Helper()
	requireNil(t, value, fmt.Sprintf("Expected nil, got: %v", value))
}

// Nilx is Nil with custom error message
func Nilx(t testing.TB, value any, format string, args ...any) {
	t.Helper()
	requireNil(t, value, fmt.Sprintf(format, args...))
}

func requireNil(t testing.TB, value any, message string) {
	if !isNil(value) {
		t.Fatal(message)
	}
}

// NotNil is a helper for ensuring that value is not nil
func NotNil(t testing.TB, value any) {
	t.Helper()
	notNil(t, value, "Expected not nil, got nil")
}

// NotNilx is NotNil with custom error message
func NotNilx(t testing.TB, value any, format string, args ...any) {
	t.Helper()
	notNil(t, value, fmt.Sprintf(format, args...))
}

func notNil(t testing.TB, value any, message string) {
	if isNil(value) {
		t.Fatal(message)
	}
}

// isNil checks if value is nil, including typed nil values
func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Func, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}


// True is a helper for ensuring that value is true
func True(t testing.TB, value bool) {
	t.Helper()
	requireTrue(t, value, "Expected true, got false")
}

// Truex is True with custom error message
func Truex(t testing.TB, value bool, format string, args ...any) {
	t.Helper()
	requireTrue(t, value, fmt.Sprintf(format, args...))
}

func requireTrue(t testing.TB, value bool, message string) {
	if !value {
		t.Fatal(message)
	}
}

// False is a helper for ensuring that value is false
func False(t testing.TB, value bool) {
	t.Helper()
	requireFalse(t, value, "Expected false, got true")
}

// Falsex is False with custom error message
func Falsex(t testing.TB, value bool, format string, args ...any) {
	t.Helper()
	requireFalse(t, value, fmt.Sprintf(format, args...))
}

func requireFalse(t testing.TB, value bool, message string) {
	if value {
		t.Fatal(message)
	}
}

// Panics is a helper for ensuring that function panics
func Panics(t testing.TB, f func()) {
	t.Helper()
	requirePanics(t, f, "Expected panic, got no panic")
}

// Panicsx is Panics with custom error message
func Panicsx(t testing.TB, f func(), format string, args ...any) {
	t.Helper()
	requirePanics(t, f, fmt.Sprintf(format, args...))
}

func requirePanics(t testing.TB, f func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal(message)
		}
	}()
	f()
}

// NotPanics is a helper for ensuring that function does not panic
func NotPanics(t testing.TB, f func()) {
	t.Helper()
	notPanics(t, f, "Expected no panic, got panic: %v")
}

// NotPanicsx is NotPanics with custom error message
func NotPanicsx(t testing.TB, f func(), format string, args ...any) {
	t.Helper()
	notPanics(t, f, fmt.Sprintf(format, args...))
}

func notPanics(t testing.TB, f func(), messageFormat string) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf(messageFormat, r)
		}
	}()
	f()
}

