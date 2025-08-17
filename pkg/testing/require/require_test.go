package require_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Equal(mt, 1, 1)
		require.Equal(mt, "hello", "hello")
		require.Equal(mt, true, true)
		
		if mt.failed {
			t.Error("Expected Equal to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.Equal(mt, 1, 2)
		
		if !mt.failed {
			t.Error("Expected Equal to fail")
		}
	})
}

func TestEqualx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Equalx(mt, 1, 1, "custom message")
		
		if mt.failed {
			t.Error("Expected Equalx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Equalx(mt, 1, 2, format, arg)
		
		if !mt.failed {
			t.Error("Expected Equalx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestEqualValues(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.EqualValues(mt, []int{1, 2, 3}, []int{1, 2, 3})
		require.EqualValues(mt, map[string]int{"a": 1}, map[string]int{"a": 1})
		
		if mt.failed {
			t.Error("Expected EqualValues to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.EqualValues(mt, []int{1, 2, 3}, []int{1, 2, 4})
		
		if !mt.failed {
			t.Error("Expected EqualValues to fail")
		}
	})
}

func TestEqualValuesx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.EqualValuesx(mt, []int{1, 2, 3}, []int{1, 2, 3}, "custom message")
		
		if mt.failed {
			t.Error("Expected EqualValuesx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.EqualValuesx(mt, []int{1, 2, 3}, []int{1, 2, 4}, format, arg)
		
		if !mt.failed {
			t.Error("Expected EqualValuesx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEqual(mt, 1, 2)
		require.NotEqual(mt, "hello", "world")
		
		if mt.failed {
			t.Error("Expected NotEqual to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEqual(mt, 1, 1)
		
		if !mt.failed {
			t.Error("Expected NotEqual to fail")
		}
	})
}

func TestNotEqualx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEqualx(mt, 1, 2, "custom message")
		
		if mt.failed {
			t.Error("Expected NotEqualx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.NotEqualx(mt, 1, 1, format, arg)
		
		if !mt.failed {
			t.Error("Expected NotEqualx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNotEqualValues(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEqualValues(mt, []int{1, 2, 3}, []int{1, 2, 4})
		require.NotEqualValues(mt, map[string]int{"a": 1}, map[string]int{"a": 2})
		
		if mt.failed {
			t.Error("Expected NotEqualValues to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEqualValues(mt, []int{1, 2, 3}, []int{1, 2, 3})
		
		if !mt.failed {
			t.Error("Expected NotEqualValues to fail")
		}
	})
}

func TestNotEqualValuesx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEqualValuesx(mt, []int{1, 2, 3}, []int{1, 2, 4}, "custom message")
		
		if mt.failed {
			t.Error("Expected NotEqualValuesx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.NotEqualValuesx(mt, []int{1, 2, 3}, []int{1, 2, 3}, format, arg)
		
		if !mt.failed {
			t.Error("Expected NotEqualValuesx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Contains(mt, "hello world", "world")
		require.Contains(mt, "test", "te")
		
		if mt.failed {
			t.Error("Expected Contains to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.Contains(mt, "hello", "world")
		
		if !mt.failed {
			t.Error("Expected Contains to fail")
		}
	})
}

func TestContainsx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Containsx(mt, "hello world", "world", "custom message")
		
		if mt.failed {
			t.Error("Expected Containsx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Containsx(mt, "hello", "world", format, arg)
		
		if !mt.failed {
			t.Error("Expected Containsx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNotContains(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotContains(mt, "hello", "world")
		require.NotContains(mt, "test", "xyz")
		
		if mt.failed {
			t.Error("Expected NotContains to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotContains(mt, "hello world", "world")
		
		if !mt.failed {
			t.Error("Expected NotContains to fail")
		}
	})
}

func TestNotContainsx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotContainsx(mt, "hello", "world", "custom message")
		
		if mt.failed {
			t.Error("Expected NotContainsx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.NotContainsx(mt, "hello world", "world", format, arg)
		
		if !mt.failed {
			t.Error("Expected NotContainsx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestEmpty(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Empty(mt, "")
		require.Empty(mt, []int{})
		require.Empty(mt, map[string]int{})
		require.Empty(mt, 0)
		require.Empty(mt, false)
		var nilPtr *int
		require.Empty(mt, nilPtr)
		
		if mt.failed {
			t.Error("Expected Empty to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.Empty(mt, "hello")
		
		if !mt.failed {
			t.Error("Expected Empty to fail")
		}
	})
}

func TestEmptyx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Emptyx(mt, "", "custom message")
		
		if mt.failed {
			t.Error("Expected Emptyx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Emptyx(mt, "hello", format, arg)
		
		if !mt.failed {
			t.Error("Expected Emptyx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNotEmpty(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEmpty(mt, "hello")
		require.NotEmpty(mt, []int{1})
		require.NotEmpty(mt, map[string]int{"a": 1})
		require.NotEmpty(mt, 1)
		require.NotEmpty(mt, true)
		
		if mt.failed {
			t.Error("Expected NotEmpty to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEmpty(mt, "")
		
		if !mt.failed {
			t.Error("Expected NotEmpty to fail")
		}
	})
}

func TestNotEmptyx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotEmptyx(mt, "hello", "custom message")
		
		if mt.failed {
			t.Error("Expected NotEmptyx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.NotEmptyx(mt, "", format, arg)
		
		if !mt.failed {
			t.Error("Expected NotEmptyx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Error(mt, errors.New("test error"))
		
		if mt.failed {
			t.Error("Expected Error to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.Error(mt, nil)
		
		if !mt.failed {
			t.Error("Expected Error to fail")
		}
	})
}

func TestErrorx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Errorx(mt, errors.New("test error"), "custom message")
		
		if mt.failed {
			t.Error("Expected Errorx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Errorx(mt, nil, format, arg)
		
		if !mt.failed {
			t.Error("Expected Errorx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestErrorIs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		baseErr := errs.New("base")
		wrappedErr := errs.Wrap(baseErr, "wrapped")
		require.ErrorIs(mt, wrappedErr, baseErr)
		
		if mt.failed {
			t.Error("Expected ErrorIs to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		err1 := errs.New("error1")
		err2 := errs.New("error2")
		require.ErrorIs(mt, err1, err2)
		
		if !mt.failed {
			t.Error("Expected ErrorIs to fail")
		}
	})
}

func TestErrorIsx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		baseErr := errs.New("base")
		wrappedErr := errs.Wrap(baseErr, "wrapped")
		require.ErrorIsx(mt, wrappedErr, baseErr, "custom message")
		
		if mt.failed {
			t.Error("Expected ErrorIsx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		err1 := errs.New("error1")
		err2 := errs.New("error2")
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.ErrorIsx(mt, err1, err2, format, arg)
		
		if !mt.failed {
			t.Error("Expected ErrorIsx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNoError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NoError(mt, nil)
		
		if mt.failed {
			t.Error("Expected NoError to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.NoError(mt, errors.New("test error"))
		
		if !mt.failed {
			t.Error("Expected NoError to fail")
		}
	})
}

func TestNoErrorx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NoErrorx(mt, nil, "custom message")
		
		if mt.failed {
			t.Error("Expected NoErrorx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.NoErrorx(mt, errors.New("test error"), format, arg)
		
		if !mt.failed {
			t.Error("Expected NoErrorx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Nil(mt, nil)
		var nilPtr *int
		require.Nil(mt, nilPtr)
		var nilSlice []int
		require.Nil(mt, nilSlice)
		
		if mt.failed {
			t.Error("Expected Nil to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.Nil(mt, "not nil")
		
		if !mt.failed {
			t.Error("Expected Nil to fail")
		}
	})
}

func TestNilx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Nilx(mt, nil, "custom message")
		
		if mt.failed {
			t.Error("Expected Nilx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Nilx(mt, "not nil", format, arg)
		
		if !mt.failed {
			t.Error("Expected Nilx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNotNil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotNil(mt, "not nil")
		require.NotNil(mt, []int{})
		require.NotNil(mt, map[string]int{})
		
		if mt.failed {
			t.Error("Expected NotNil to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotNil(mt, nil)
		
		if !mt.failed {
			t.Error("Expected NotNil to fail")
		}
	})
}

func TestNotNilx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotNilx(mt, "not nil", "custom message")
		
		if mt.failed {
			t.Error("Expected NotNilx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.NotNilx(mt, nil, format, arg)
		
		if !mt.failed {
			t.Error("Expected NotNilx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestTrue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.True(mt, true)
		
		if mt.failed {
			t.Error("Expected True to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.True(mt, false)
		
		if !mt.failed {
			t.Error("Expected True to fail")
		}
	})
}

func TestTruex(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Truex(mt, true, "custom message")
		
		if mt.failed {
			t.Error("Expected Truex to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Truex(mt, false, format, arg)
		
		if !mt.failed {
			t.Error("Expected Truex to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestFalse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.False(mt, false)
		
		if mt.failed {
			t.Error("Expected False to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.False(mt, true)
		
		if !mt.failed {
			t.Error("Expected False to fail")
		}
	})
}

func TestFalsex(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Falsex(mt, false, "custom message")
		
		if mt.failed {
			t.Error("Expected Falsex to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Falsex(mt, true, format, arg)
		
		if !mt.failed {
			t.Error("Expected Falsex to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestPanics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Panics(mt, func() {
			panic("test panic")
		})
		
		if mt.failed {
			t.Error("Expected Panics to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.Panics(mt, func() {
			// no panic
		})
		
		if !mt.failed {
			t.Error("Expected Panics to fail")
		}
	})
}

func TestPanicsx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.Panicsx(mt, func() {
			panic("test panic")
		}, "custom message")
		
		if mt.failed {
			t.Error("Expected Panicsx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.Panicsx(mt, func() {
			// no panic
		}, format, arg)
		
		if !mt.failed {
			t.Error("Expected Panicsx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

func TestNotPanics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotPanics(mt, func() {
			// no panic
		})
		
		if mt.failed {
			t.Error("Expected NotPanics to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotPanics(mt, func() {
			panic("test panic")
		})
		
		if !mt.failed {
			t.Error("Expected NotPanics to fail")
		}
	})
}

func TestNotPanicsx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mt := &mockTesting{}
		require.NotPanicsx(mt, func() {
			// no panic
		}, "custom message")
		
		if mt.failed {
			t.Error("Expected NotPanicsx to pass")
		}
	})

	t.Run("failure", func(t *testing.T) {
		mt := &mockTesting{}
		format := "custom message %s"
		arg := "test"
		expected := fmt.Sprintf(format, arg)
		require.NotPanicsx(mt, func() {
			panic("test panic")
		}, format, arg)
		
		if !mt.failed {
			t.Error("Expected NotPanicsx to fail")
		}
		if len(mt.logs) == 0 || !strings.Contains(mt.logs[0], expected) {
			t.Errorf("Expected log to contain formatted message %q", expected)
		}
	})
}

// mockTesting implements testing.TB for testing purposes
type mockTesting struct {
	testing.TB
	failed bool
	logs   []string
}

func (m *mockTesting) Cleanup(func()) {}
func (m *mockTesting) Error(args ...any) {
	m.failed = true
}
func (m *mockTesting) Errorf(format string, args ...any) {
	m.failed = true
}
func (m *mockTesting) Fail() {
	m.failed = true
}
func (m *mockTesting) FailNow() {
	m.failed = true
}
func (m *mockTesting) Failed() bool {
	return m.failed
}
func (m *mockTesting) Fatal(args ...any) {
	m.failed = true
}
func (m *mockTesting) Fatalf(format string, args ...any) {
	m.failed = true
}
func (m *mockTesting) Helper() {}
func (m *mockTesting) Log(args ...any) {
	m.logs = append(m.logs, args[0].(string))
}
func (m *mockTesting) Logf(format string, args ...any) {
	m.logs = append(m.logs, format)
}
func (m *mockTesting) Name() string {
	return "mock"
}
func (m *mockTesting) Setenv(key, value string) {}
func (m *mockTesting) Skip(args ...any) {}
func (m *mockTesting) SkipNow() {}
func (m *mockTesting) Skipf(format string, args ...any) {}
func (m *mockTesting) Skipped() bool {
	return false
}
func (m *mockTesting) TempDir() string {
	return "/tmp"
}

func TestIsEmpty(t *testing.T) {
	// Test cases for different empty values
	emptyValues := []any{
		"",
		[]int{},
		[]string{},
		map[string]int{},
		make(chan int),
		false,
		0,
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		uintptr(0),
		float32(0),
		float64(0),
		complex64(0),
		complex128(0),
		(*int)(nil),
		([]int)(nil),
		(map[string]int)(nil),
		(chan int)(nil),
		(func())(nil),
		struct{}{},
	}

	for _, val := range emptyValues {
		mt := &mockTesting{}
		require.Empty(mt, val)
		if mt.failed {
			t.Errorf("Expected value %v to be empty", val)
		}
	}

	// Test cases for non-empty values
	nonEmptyValues := []any{
		"hello",
		[]int{1},
		map[string]int{"a": 1},
		true,
		1,
		-1,
		1.0,
		complex(1, 0),
		&struct{}{},
		struct{ A int }{A: 1},
	}

	for _, val := range nonEmptyValues {
		mt := &mockTesting{}
		require.NotEmpty(mt, val)
		if mt.failed {
			t.Errorf("Expected value %v to be non-empty", val)
		}
	}
}

func TestIsNil(t *testing.T) {
	// Test cases for nil values
	nilValues := []any{
		nil,
		(*int)(nil),
		([]int)(nil),
		(map[string]int)(nil),
		(chan int)(nil),
		(func())(nil),
		(any)(nil),
	}

	for _, val := range nilValues {
		mt := &mockTesting{}
		require.Nil(mt, val)
		if mt.failed {
			t.Errorf("Expected value %v to be nil", val)
		}
	}

	// Test cases for non-nil values
	nonNilValues := []any{
		"",
		0,
		false,
		[]int{},
		map[string]int{},
		&struct{}{},
		struct{}{},
	}

	for _, val := range nonNilValues {
		mt := &mockTesting{}
		require.NotNil(mt, val)
		if mt.failed {
			t.Errorf("Expected value %v to be non-nil", val)
		}
	}
}
