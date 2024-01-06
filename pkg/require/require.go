package require

import "testing"

// Equal is a helper for comparing to variables in non deep-equal way
func Equal[T comparable](t *testing.T, expect, got T) {
	if expect != got {
		logDiff(t, expect, got)
		t.Fail()
	}
}

// EqualWithMessage is a helper for comparing to variables in non deep-equal way.
// Will print provided message
func EqualWithMessage[T comparable](t *testing.T, expect, got T, msg string, args ...any) {
	if expect != got {
		logDiff(t, expect, got)
		t.Fatalf(msg, args...)
	}
}

func logDiff[T comparable](t *testing.T, expected, got T) {
	t.Logf("\nExpected: %v\nGot:      %v", expected, got)
}
