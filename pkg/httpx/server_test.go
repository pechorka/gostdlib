package httpx_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pechorka/gostdlib/pkg/httpx"
	"github.com/pechorka/gostdlib/pkg/testing/require"
)

type testData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestReadJSON(t *testing.T) {
	t.Run("successful parse", func(t *testing.T) {
		expected := testData{
			Message: "success",
			Code:    200,
		}
		jsonData, err := json.Marshal(expected)
		require.NoError(t, err)

		reader := bytes.NewReader(jsonData)
		result, err := httpx.ReadJSON[testData](reader)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		reader := strings.NewReader(`{"invalid json`)
		_, err := httpx.ReadJSON[testData](reader)
		require.Error(t, err)
	})

	t.Run("empty body", func(t *testing.T) {
		reader := strings.NewReader("")
		_, err := httpx.ReadJSON[testData](reader)
		require.Error(t, err)
	})

	t.Run("nil reader", func(t *testing.T) {
		_, err := httpx.ReadJSON[testData](nil)
		require.Error(t, err)
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("successful write", func(t *testing.T) {
		data := testData{
			Message: "success",
			Code:    200,
		}

		w := httptest.NewRecorder()
		err := httpx.WriteJSON(w, data)
		require.NoError(t, err)

		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var result testData
		err = json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, data, result)
	})

	t.Run("write nil value", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := httpx.WriteJSON[any](w, nil)
		require.NoError(t, err)

		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.Equal(t, "null\n", w.Body.String())
	})

	t.Run("write to failed response writer", func(t *testing.T) {
		// Create a custom ResponseWriter that fails on write
		failWriter := &failingResponseWriter{
			header: make(http.Header),
		}

		data := testData{
			Message: "success",
			Code:    200,
		}

		err := httpx.WriteJSON(failWriter, data)
		require.Error(t, err)
	})
}

// failingResponseWriter is a custom ResponseWriter that always fails to write
type failingResponseWriter struct {
	header http.Header
}

func (f *failingResponseWriter) Header() http.Header {
	return f.header
}

func (f *failingResponseWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func (f *failingResponseWriter) WriteHeader(statusCode int) {
	// Do nothing
}
