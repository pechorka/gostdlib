package sqlx

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestKVDriver(t *testing.T) {
	// Open connection
	db, err := sql.Open("kv", "")
	require.NoError(t, err)
	defer db.Close()

	t.Run("SET and GET operations", func(t *testing.T) {
		// Test setting a value
		result, err := db.Exec("SET ? = ?", "test_key", "test_value")
		require.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAffected)

		// Test getting the value
		var value string
		err = db.QueryRow("GET ?", "test_key").Scan(&value)
		require.NoError(t, err)
		require.Equal(t, "test_value", value)
	})

	t.Run("GET non-existent key", func(t *testing.T) {
		var value string
		err := db.QueryRow("GET ?", "non_existent_key").Scan(&value)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("SET validation", func(t *testing.T) {
		// Test with wrong number of arguments
		_, err := db.Exec("SET ?", "only_key")
		require.Error(t, err)

		// Test with non-string key
		_, err = db.Exec("SET ? = ?", 123, "value")
		require.Error(t, err)

		// Test with non-string value
		_, err = db.Exec("SET ? = ?", "key", 123)
		require.Error(t, err)
	})

	t.Run("GET validation", func(t *testing.T) {
		// Test with wrong number of arguments
		err := db.QueryRow("GET ?, ?", "key1", "key2").Scan(new(string))
		require.Error(t, err)

		// Test with non-string key
		err = db.QueryRow("GET ?", 123).Scan(new(string))
		require.Error(t, err)
	})

	t.Run("Unsupported operations", func(t *testing.T) {
		// Test unsupported query operation
		_, err := db.Exec("DELETE ?", "key")
		require.Error(t, err)

		// Test unsupported query operation
		err = db.QueryRow("SELECT ?", "key").Scan(new(string))
		require.Error(t, err)
	})

	t.Run("Multiple operations", func(t *testing.T) {
		// Set multiple key-value pairs
		pairs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		for k, v := range pairs {
			_, err := db.Exec("SET ? = ?", k, v)
			require.NoError(t, err)
		}

		// Verify all values
		for k, expectedV := range pairs {
			var actualV string
			err := db.QueryRow("GET ?", k).Scan(&actualV)
			require.NoError(t, err)
			require.Equal(t, expectedV, actualV)
		}
	})

	t.Run("Concurrent operations", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func(i int) {
				key := fmt.Sprintf("concurrent_key_%d", i)
				value := fmt.Sprintf("concurrent_value_%d", i)

				_, err := db.Exec("SET ? = ?", key, value)
				require.NoError(t, err)

				var readValue string
				err = db.QueryRow("GET ?", key).Scan(&readValue)
				require.NoError(t, err)
				require.Equal(t, value, readValue)

				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}
