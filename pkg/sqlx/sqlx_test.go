package sqlx

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func Test_mapDestValues(t *testing.T) {
	t.Run("all columns found", func(t *testing.T) {
		var dest struct {
			ID     int    `db:"id"`
			Name   string `db:"name"`
			Banned bool   `db:"banned"`
			hidden string
		}

		cols := []string{"id", "name", "banned"}
		columnIndex, err := columnIndexToDestFieldIndex(cols, reflect.TypeOf(dest))
		require.NoError(t, err)

		destValues, err := mapDestValue(reflect.ValueOf(&dest).Elem(), columnIndex)
		require.NoError(t, err)
		require.EqualValues(t, &dest.ID, destValues[0])
		require.EqualValues(t, &dest.Name, destValues[1])
	})
	t.Run("nil columns", func(t *testing.T) {
		var dest struct {
			ID       int            `db:"id"`
			Nickname *string        `db:"nickname"`
			Map      map[string]int `db:"map"`
		}
		cols := []string{"id", "nickname", "map"}
		columnIndex, err := columnIndexToDestFieldIndex(cols, reflect.TypeOf(dest))
		require.NoError(t, err)

		destValues, err := mapDestValue(reflect.ValueOf(&dest).Elem(), columnIndex)
		require.NoError(t, err)
		require.EqualValues(t, &dest.ID, destValues[0])
		require.EqualValues(t, &dest.Nickname, destValues[1])
		require.NotNil(t, dest.Nickname)
		require.EqualValues(t, &dest.Map, destValues[2])
		require.NotNil(t, dest.Map)
	})

	t.Run("non struct dest", func(t *testing.T) {
		var dest string
		cols := []string{"id"}
		columnIndex, err := columnIndexToDestFieldIndex(cols, reflect.TypeOf(dest))
		require.NoError(t, err)

		destValues, err := mapDestValue(reflect.ValueOf(&dest).Elem(), columnIndex)
		require.NoError(t, err)
		require.EqualValues(t, &dest, destValues[0])
	})
}

func Test_Get(t *testing.T) {
	ctx, db := newTestCtx(t)

	var dest struct {
		Text      string    `db:"text"`
		Num       int       `db:"num"`
		Bool      bool      `db:"bool"`
		Scannable scannable `db:"scannable"`
		Date      time.Time `db:"date"`
	}

	now := time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO test(text, num, bool, scannable, date) VALUES(?, ?, ?, ?, ?)
	`, "text", 1, true, scannable{"string"}, now)
	require.NoError(t, err)

	err = db.Get(&dest, `SELECT * FROM test WHERE text = ?`, "text")
	require.NoError(t, err)
	require.EqualValues(t, "text", dest.Text)
	require.EqualValues(t, 1, dest.Num)
	require.EqualValues(t, true, dest.Bool)
	require.EqualValues(t, scannable{"string"}, dest.Scannable)
	require.EqualValues(t, now, dest.Date)
}

type scannable struct {
	String string `json:"string"`
}

func (s *scannable) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", src)
	}
	return json.Unmarshal(b, s)
}

func (s *scannable) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type testCtx struct {
	context.Context
}

func newTestCtx(t *testing.T) (*testCtx, *DB) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)
	db, err := Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE test(
			text TEXT NOT NULL,
			num INTEGER NOT NULL,
			bool BOOLEAN NOT NULL,
			scannable JSONB NOT NULL,
			date TIMESTAMP NOT NULL
		)	
	`)
	require.NoError(t, err)

	return &testCtx{ctx}, db
}
