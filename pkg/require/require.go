package require

import (
	"testing"
)

// Equal is a helper for comparing 2 variables in non deep-equal way
func Equal[T comparable](t *testing.T, expect, got T) {
	if expect != got {
		logDiff(t, expect, got)
		t.Fail()
	}
}

// EqualWithMessage is a helper for comparing 2 variables in non deep-equal way.
// Will print provided message
func EqualWithMessage[T comparable](t *testing.T, expect, got T, msg string, args ...any) {
	if expect != got {
		logDiff(t, expect, got)
		t.Fatalf(msg, args...)
	}
}

// Error is a helper for ensuring that error is not nil
func Error(t *testing.T, err error) {
	if err == nil {
		t.Fatal("\nExpected error")
	}
}

// NoError is a helper for ensuring that error is nil
func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("\nExpected no error, got: %v", err)
	}
}

func logDiff[T comparable](t *testing.T, expected, got T) {
	t.Logf("\nExpected: %v\nGot:      %v", expected, got)
}
