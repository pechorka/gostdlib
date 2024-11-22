package require

import (
	"reflect"
	"testing"

	"github.com/pechorka/gostdlib/pkg/testing/differ"
)

// Equal is a helper for comparing 2 variables in non deep-equal way
func Equal[T comparable](t *testing.T, expect, got T) {
	if expect != got {
		t.Log(differ.Diff(expect, got))
		t.Fail()
	}
}

// DeepEqual compares two values using reflect.DeepEqual and shows detailed differences
func DeepEqual(t *testing.T, expect, got interface{}) {
	if !reflect.DeepEqual(expect, got) {
		t.Log(differ.Diff(expect, got))
		t.Fail()
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

// Nil is a helper for ensuring that value is nil
func Nil(t *testing.T, value any) {
	if value != nil {
		t.Fatalf("\nExpected nil, got: %v", value)
	}
}
