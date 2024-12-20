package require

import (
	"reflect"
	"testing"

	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/testing/differ"
)

// Equal is a helper for comparing 2 variables in non deep-equal way
func Equal[T comparable](t *testing.T, expect, got T) {
	t.Helper()
	if expect != got {
		t.Log(differ.Diff(expect, got))
		t.Fail()
	}
}

// EqualValues compares two values that can't be compared by == operator
func EqualValues(t *testing.T, expect, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expect, got) {
		t.Log(differ.Diff(expect, got))
		t.Fail()
	}
}

// Error is a helper for ensuring that error is not nil
func Error(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("\nExpected error")
	}
}

// ErrorIs is a helper for ensuring that error is target error
func ErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errs.Is(err, target) {
		t.Fatalf("\nExpected error %v, got %v", target, err)
	}
}

// NoError is a helper for ensuring that error is nil
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("\nExpected no error, got: %v", err)
	}
}

// Nil is a helper for ensuring that value is nil
func Nil(t *testing.T, value any) {
	t.Helper()
	if value != nil {
		t.Fatalf("\nExpected nil, got: %v", value)
	}
}

func True(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Fatal("\nExpected true, got false")
	}
}

func False(t *testing.T, value bool) {
	t.Helper()
	if value {
		t.Fatal("\nExpected false, got true")
	}
}

func Panics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("\nExpected panic, got no panic")
		}
	}()
	f()
}

func NoPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("\nExpected no panic, got panic: %v", r)
		}
	}()
	f()
}
