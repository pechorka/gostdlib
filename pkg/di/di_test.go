package di

import (
	"errors"
	"strings"
	"testing"

	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestNewProvider(t *testing.T) {
	t.Run("non-function provider", func(t *testing.T) {
		_, err := NewProvider("not a function")
		require.Error(t, err)
		require.Equal(t, "failed to parse providers: 0th provider is not a function, got string", err.Error())
	})

	t.Run("provider with no output", func(t *testing.T) {
		_, err := NewProvider(func() {})
		require.Error(t, err)
		require.Equal(t, "failed to parse providers: 0th provider 1 has no output", err.Error())
	})

	t.Run("provider with more than two outputs", func(t *testing.T) {
		providerFunc := func() (string, string, int) { return "", "", 0 }
		_, err := NewProvider(providerFunc)
		require.Error(t, err)
		require.Equal(t, "failed to parse providers: 0th provider 1 has more than two outputs. Provider must return a single value or a value and an error", err.Error())
	})

	t.Run("provider with wrong second output type", func(t *testing.T) {
		providerFunc := func() (string, string) { return "", "" }
		_, err := NewProvider(providerFunc)
		require.Error(t, err)
		require.Equal(t, "failed to parse providers: 0th provider 1 has two outputs, but the second one is not an error", err.Error())
	})

	t.Run("duplicate provider", func(t *testing.T) {
		providerFunc1 := func() string { return "" }
		providerFunc2 := func() string { return "" }
		_, err := NewProvider(providerFunc1, providerFunc2)
		require.Error(t, err)
		require.Equal(t, "failed to parse providers: 1th provider 2 returns the same type string as provider 1", err.Error())
	})
	t.Run("missing dependency", func(t *testing.T) {
		providerFunc := func(int) string { return "" }
		_, err := NewProvider(providerFunc)
		require.Error(t, err)
		require.Equal(t, "all deps must be provided: dependency int is not provided", err.Error())
	})
	t.Run("cyclic dependency", func(t *testing.T) {
		_, err := NewProvider(
			func(i int) string { return "" },
			func(s string) int { return 0 },
		)
		require.Error(t, err)
		expectedPrefix := "should not have cyclic dependencies"
		require.True(t, strings.HasPrefix(err.Error(), expectedPrefix))
	})
	t.Run("valid providers", func(t *testing.T) {
		providerFunc1 := func() string { return "hello" }
		providerFunc2 := func(s string) int { return len(s) }
		_, err := NewProvider(providerFunc1, providerFunc2)
		require.NoError(t, err)
	})
}

func TestProvider_Provide(t *testing.T) {
	t.Run("non-pointer destination", func(t *testing.T) {
		provider, err := NewProvider(func() string { return "hello" })
		require.NoError(t, err)

		err = provider.Provide(struct{}{})
		require.Error(t, err)
	})

	t.Run("non-struct pointer destination", func(t *testing.T) {
		provider, err := NewProvider(func() string { return "hello" })
		require.NoError(t, err)

		err = provider.Provide(new(string))
		require.Error(t, err)
	})

	t.Run("provider returns error", func(t *testing.T) {
		errText := "provider error"
		provider, err := NewProvider(func() (string, error) {
			return "", errors.New(errText)
		})
		require.NoError(t, err)

		err = provider.Provide(&struct{ S string }{})
		require.Error(t, err)
		contains := strings.Contains(err.Error(), errText)
		require.True(t, contains)
	})
	t.Run("provider dependency returns error", func(t *testing.T) {
		errText := "provider dependency error"
		depProvider := func() (string, error) {
			return "", errors.New(errText)
		}
		mainProvider := func(s string) int {
			return len(s)
		}
		provider, err := NewProvider(depProvider, mainProvider)
		require.NoError(t, err)

		err = provider.Provide(&struct {
			I int
		}{})
		require.Error(t, err)
		contains := strings.Contains(err.Error(), errText)
		require.True(t, contains)
	})

	t.Run("missing field provider", func(t *testing.T) {
		provider, err := NewProvider(func() string { return "hello" })
		require.NoError(t, err)

		err = provider.Provide(&struct {
			S string
			I int
		}{})
		require.Error(t, err)
		require.Equal(t, "failed to resolve field I: no provider found for type int", err.Error())
	})

	t.Run("successful injection", func(t *testing.T) {
		provider, err := NewProvider(
			func() string { return "hello" },
			func(s string) int { return len(s) },
		)
		require.NoError(t, err)

		dst := &struct {
			S string
			I int
		}{}
		err = provider.Provide(dst)
		require.NoError(t, err)

		require.Equal(t, "hello", dst.S)
		require.Equal(t, 5, dst.I)
	})
}

func TestComplexDependencies(t *testing.T) {
	type Config struct {
		DBHost string
		DBPort int
	}
	type Database struct {
		Config Config
		URL    string
	}
	type Service struct {
		DB     Database
		Active bool
	}

	provider, err := NewProvider(
		func() Config {
			return Config{
				DBHost: "localhost",
				DBPort: 5432,
			}
		},
		func(cfg Config) Database {
			return Database{
				Config: cfg,
				URL:    "postgresql://" + cfg.DBHost,
			}
		},
		func(db Database) Service {
			return Service{
				DB:     db,
				Active: true,
			}
		},
	)
	require.NoError(t, err)

	dst := &struct {
		Cfg Config
		DB  Database
		Svc Service
	}{}

	err = provider.Provide(dst)
	require.NoError(t, err)

	require.Equal(t, "localhost", dst.Cfg.DBHost)
	require.Equal(t, "postgresql://localhost", dst.DB.URL)
	require.True(t, dst.Svc.Active)
}
