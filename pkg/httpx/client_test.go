package httpx_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pechorka/gostdlib/pkg/httpx"
	"github.com/pechorka/gostdlib/pkg/testing/require"
)

type testResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestGetJSON(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		expected := testResponse{
			Message: "success",
			Code:    200,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)

			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(expected)
			require.NoError(t, err)
		}))
		defer server.Close()

		resp, err := httpx.GetJSON[testResponse](context.Background(), server.URL)
		require.NoError(t, err)
		require.Equal(t, expected, resp)
	})

	t.Run("invalid URL", func(t *testing.T) {
		_, err := httpx.GetJSON[testResponse](context.Background(), "http://invalid-url")
		require.Error(t, err)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		_, err := httpx.GetJSON[testResponse](context.Background(), server.URL)
		require.Error(t, err)
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"invalid json`))
		}))
		defer server.Close()

		_, err := httpx.GetJSON[testResponse](context.Background(), server.URL)
		require.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}))
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := httpx.GetJSON[testResponse](ctx, server.URL)
		require.ErrorIs(t, err, context.Canceled)
	})
}
