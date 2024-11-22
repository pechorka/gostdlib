package errs_test

import (
	"testing"

	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestWrap(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		baseErr := errs.New("base")
		wrapped := errs.Wrap(baseErr, "wrap")

		require.Equal(t, baseErr, errs.Unwrap(wrapped))
		require.Equal(t, true, errs.Is(wrapped, baseErr))
	})
	t.Run("nil error", func(t *testing.T) {
		require.Nil(t, errs.Wrap(nil, "wrap"))
	})
}

func TestWrapf(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		baseErr := errs.New("base")
		wrapped := errs.Wrapf(baseErr, "wrap")

		require.Equal(t, baseErr, errs.Unwrap(wrapped))
		require.Equal(t, true, errs.Is(wrapped, baseErr))
	})

	t.Run("args", func(t *testing.T) {
		baseErr := errs.New("base error")
		arg := "argName"
		wrapped := errs.Wrapf(baseErr, "wrap with arg %s", arg)

		require.Equal(t, baseErr, errs.Unwrap(wrapped))
		require.Equal(t, true, errs.Is(wrapped, baseErr))
		require.Equal(t, "wrap with arg argName: base error", wrapped.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		require.Nil(t, errs.Wrapf(nil, "wrap"))
	})
}
