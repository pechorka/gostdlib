package uuid

import (
	"testing"
	"time"

	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestUUID(t *testing.T) {
	t.Run("String format", func(t *testing.T) {
		uuid := UUID{
			0x12, 0x34, 0x56, 0x78,
			0x9a, 0xbc,
			0xde, 0xf0,
			0x12, 0x34,
			0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
		}
		expected := "12345678-9abc-def0-1234-56789abcdef0"
		require.Equal(t, expected, uuid.String())
	})

	t.Run("NewV4", func(t *testing.T) {
		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			uuid, err := NewV4(mockReader)
			require.NoError(t, err)

			expected := "12345678-9abc-4ef0-9234-56789abcdef0" // Note: version and variant bits modified
			require.Equal(t, expected, uuid.String())

			// Check version bits (version 4)
			require.Equal(t, byte(0x4), uuid[6]>>4)
			// Check variant bits (RFC4122)
			require.Equal(t, byte(0x2), uuid[8]>>6)
		})

		t.Run("reader error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			_, err := NewV4(errReader)
			require.Error(t, err)
		})
	})

	t.Run("MustV4", func(t *testing.T) {
		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			uuid := MustV4(mockReader)
			expected := "12345678-9abc-4ef0-9234-56789abcdef0"
			require.Equal(t, expected, uuid.String())
		})

		t.Run("panics on error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			require.Panics(t, func() {
				MustV4(errReader)
			})
		})
	})

	t.Run("NewV4String", func(t *testing.T) {
		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			uuidStr, err := NewV4String(mockReader)
			require.NoError(t, err)
			require.Equal(t, "12345678-9abc-4ef0-9234-56789abcdef0", uuidStr)
		})

		t.Run("reader error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			_, err := NewV4String(errReader)
			require.Error(t, err)
		})
	})

	t.Run("MustV4String", func(t *testing.T) {
		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			uuidStr := MustV4String(mockReader)
			require.Equal(t, "12345678-9abc-4ef0-9234-56789abcdef0", uuidStr)
		})

		t.Run("panics on error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			require.Panics(t, func() {
				MustV4String(errReader)
			})
		})
	})

	checkTimestamp := func(uuid UUID, before, after int64) {
		timestamp := int64(uuid[0])<<40 | int64(uuid[1])<<32 | int64(uuid[2])<<24 | int64(uuid[3])<<16 | int64(uuid[4])<<8 | int64(uuid[5])
		require.True(t, before <= timestamp && timestamp <= after)
	}
	t.Run("NewV7", func(t *testing.T) {

		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			before := time.Now().UnixMilli()
			uuid, err := NewV7(mockReader)
			after := time.Now().UnixMilli()
			require.NoError(t, err)

			checkTimestamp(uuid, before, after)

			// Check version bits (version 7)
			require.Equal(t, byte(0x7), uuid[6]>>4)
			// Check variant bits (RFC4122)
			require.Equal(t, byte(0x2), uuid[8]>>6)
		})

		t.Run("reader error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			_, err := NewV7(errReader)
			require.Error(t, err)
		})
	})

	t.Run("MustV7", func(t *testing.T) {
		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			before := time.Now().UnixMilli()
			uuid := MustV7(mockReader)
			after := time.Now().UnixMilli()

			checkTimestamp(uuid, before, after)

			require.Equal(t, byte(0x7), uuid[6]>>4) // Check version bits (version 7)
			require.Equal(t, byte(0x2), uuid[8]>>6) // Check variant bits (RFC4122)
		})

		t.Run("panics on error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			require.Panics(t, func() {
				MustV7(errReader)
			})
		})
	})

	t.Run("NewV7String", func(t *testing.T) {
		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			uuidStr, err := NewV7String(mockReader)
			require.NoError(t, err)
			// Verify it's a valid UUID string format
			require.Equal(t, 36, len(uuidStr))
			require.Equal(t, byte('7'), uuidStr[14])
		})

		t.Run("reader error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			_, err := NewV7String(errReader)
			require.Error(t, err)
		})
	})

	t.Run("MustV7String", func(t *testing.T) {
		t.Run("valid generation", func(t *testing.T) {
			mockReader := &mockReader{
				data: []byte{
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
					0x12, 0x34, 0x56, 0x78,
					0x9a, 0xbc, 0xde, 0xf0,
				},
			}
			uuidStr := MustV7String(mockReader)
			// Verify it's a valid UUID string format
			require.Equal(t, 36, len(uuidStr))
			require.Equal(t, byte('7'), uuidStr[14])
		})

		t.Run("panics on error", func(t *testing.T) {
			errReader := &errorReader{err: errs.New("reader error")}
			require.Panics(t, func() {
				MustV7String(errReader)
			})
		})
	})
}

// mockReader is a mock io.Reader that returns predefined data
type mockReader struct {
	data []byte
	pos  int
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	remaining := len(r.data) - r.pos
	if remaining == 0 {
		return 0, nil
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// errorReader is a mock io.Reader that always returns an error
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}
