package cryptokit

import (
	"bytes"
	"testing"

	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestAesEncrypter_Encrypt(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		key := bytes.Repeat([]byte("a"), 32)
		cryptor, err := NewEncrypter(key)
		require.NoError(t, err)

		data := []byte("373512635")
		encrypted, err := cryptor.Encrypt(data)
		require.NoError(t, err)

		decrypted, err := cryptor.Decrypt(encrypted)
		require.NoError(t, err)

		require.Equal(t, string(data), string(decrypted))
	})
}
