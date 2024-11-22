package env

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/testing/require"
)

type testConfig struct {
	Db      DbConfig       `env:"DB"`
	Kafka   KafkaConfig    `env:"KAFKA"`
	Timeout CustomDuration `env:"TIMEOUT" default:"10d"`
}

type DbConfig struct {
	Host string `env:"HOST" required:"true"`
	Port int    `env:"PORT" default:"5432"`
}

type KafkaConfig struct {
	Brokers []string `env:"BROKERS" sep:","`
}

type CustomDuration time.Duration

func (d *CustomDuration) UnmarshalEnv(v string) error {
	if v == "" {
		return nil
	}
	if strings.HasSuffix(v, "d") {
		v = strings.TrimSuffix(v, "d")
		days, err := strconv.Atoi(v)
		if err != nil {
			return errs.Wrap(err, "failed to parse duration")
		}
		*d = CustomDuration(time.Duration(days) * 24 * time.Hour)
		return nil
	}
	duration, err := time.ParseDuration(v)
	if err != nil {
		return errs.Wrap(err, "failed to parse duration")
	}
	*d = CustomDuration(duration)
	return nil
}

func TestParseConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		prepareEnv(t,
			"DB_HOST", "localhost",
			"DB_PORT", "5432",
			"KAFKA_BROKERS", "localhost:9092,localhost:9093",
			"TIMEOUT", "10d",
		)

		cfg := &testConfig{}
		err := ParseConfig(cfg)
		require.NoError(t, err)

		require.Equal(t, "localhost", cfg.Db.Host)
		require.Equal(t, 5432, cfg.Db.Port)
		require.EqualValues(t, []string{"localhost:9092", "localhost:9093"}, cfg.Kafka.Brokers)
		require.Equal(t, CustomDuration(10*24*time.Hour), cfg.Timeout)
	})

	t.Run("missing required field", func(t *testing.T) {
		prepareEnv(t,
			"DB_PORT", "5432",
		)

		cfg := &testConfig{}
		err := ParseConfig(cfg)
		require.Error(t, err)
		require.True(t, strings.HasSuffix(err.Error(), "DB_HOST is not set"))
	})

	t.Run("default values", func(t *testing.T) {
		prepareEnv(t,
			"DB_HOST", "localhost",
		)

		cfg := &testConfig{}
		err := ParseConfig(cfg)
		require.NoError(t, err)

		require.Equal(t, "localhost", cfg.Db.Host)
		require.Equal(t, 5432, cfg.Db.Port)                            // Default value
		require.Equal(t, CustomDuration(10*24*time.Hour), cfg.Timeout) // Default value
	})

	t.Run("invalid input types", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func()
			cleanup func()
		}{
			{
				name: "invalid port number",
				setup: func() {
					os.Setenv("DB_HOST", "localhost")
					os.Setenv("DB_PORT", "invalid")
				},
				cleanup: func() {
					os.Unsetenv("DB_HOST")
					os.Unsetenv("DB_PORT")
				},
			},
			{
				name: "invalid duration",
				setup: func() {
					os.Setenv("DB_HOST", "localhost")
					os.Setenv("TIMEOUT", "invalid")
				},
				cleanup: func() {
					os.Unsetenv("DB_HOST")
					os.Unsetenv("TIMEOUT")
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setup()
				defer tt.cleanup()

				cfg := &testConfig{}
				err := ParseConfig(cfg)
				require.Error(t, err)
			})
		}
	})

	t.Run("non-pointer input", func(t *testing.T) {
		cfg := testConfig{}
		err := ParseConfig(cfg)
		require.Error(t, err)
		require.Equal(t, err.Error(), "config must be a non-nil pointer")
	})

	t.Run("nil input", func(t *testing.T) {
		err := ParseConfig(nil)
		require.Error(t, err)
		require.Equal(t, err.Error(), "config must be a non-nil pointer")
	})

	t.Run("custom unmarshaler", func(t *testing.T) {
		tests := []struct {
			name     string
			envValue string
			expected CustomDuration
		}{
			{
				name:     "days format",
				envValue: "5d",
				expected: CustomDuration(5 * 24 * time.Hour),
			},
			{
				name:     "hours format",
				envValue: "48h",
				expected: CustomDuration(48 * time.Hour),
			},
			{
				name:     "minutes format",
				envValue: "90m",
				expected: CustomDuration(90 * time.Minute),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("TIMEOUT", tt.envValue)
				defer func() {
					os.Unsetenv("DB_HOST")
					os.Unsetenv("TIMEOUT")
				}()

				cfg := &testConfig{}
				err := ParseConfig(cfg)
				require.NoError(t, err)
				require.Equal(t, tt.expected, cfg.Timeout)
			})
		}
	})
}

func prepareEnv(t *testing.T, envValue ...string) {
	t.Helper()
	if len(envValue)%2 != 0 {
		t.Fatalf("envValue must be a pair of key and value")
	}
	for i := 0; i < len(envValue); i += 2 {
		os.Setenv(envValue[i], envValue[i+1])
		t.Cleanup(func() {
			os.Unsetenv(envValue[i])
		})
	}
}
