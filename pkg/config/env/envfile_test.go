package env

import (
	"os"
	"testing"

	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func Test_exportDotEnv(t *testing.T) {
	t.Run("valid env vars", func(t *testing.T) {
		input := []byte(`FOO=bar
BAZ=qux`)
		err := exportDotEnv(input)
		require.NoError(t, err)

		require.Equal(t, "bar", os.Getenv("FOO"))
		require.Equal(t, "qux", os.Getenv("BAZ"))

		// Cleanup
		os.Unsetenv("FOO")
		os.Unsetenv("BAZ")
	})

	t.Run("empty lines and comments are skipped", func(t *testing.T) {
		input := []byte(`FOO=bar

# this is a comment
BAZ=qux`)
		err := exportDotEnv(input)
		require.NoError(t, err)

		require.Equal(t, "bar", os.Getenv("FOO"))
		require.Equal(t, "qux", os.Getenv("BAZ"))

		// Cleanup
		os.Unsetenv("FOO")
		os.Unsetenv("BAZ")
	})

	t.Run("invalid line format returns error", func(t *testing.T) {
		input := []byte(`FOO=bar
INVALID_LINE
BAZ=qux`)
		err := exportDotEnv(input)
		require.Error(t, err)
	})

	t.Run("empty value is allowed", func(t *testing.T) {
		input := []byte(`EMPTY=`)
		err := exportDotEnv(input)
		require.NoError(t, err)

		require.Equal(t, "", os.Getenv("EMPTY"))

		// Cleanup
		os.Unsetenv("EMPTY")
	})

	t.Run("value can contain equals sign", func(t *testing.T) {
		input := []byte(`URL=http://example.com?foo=bar`)
		err := exportDotEnv(input)
		require.NoError(t, err)

		require.Equal(t, "http://example.com?foo=bar", os.Getenv("URL"))

		// Cleanup
		os.Unsetenv("URL")
	})
}
