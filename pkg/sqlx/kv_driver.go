package sqlx

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

func init() {
	sql.Register("kv", &KVDriver{})
}

var (
	kvStorage map[string]string = make(map[string]string, 10)
	mu        sync.RWMutex
)

// KVDriver implements database/sql/driver.Driver interface
type KVDriver struct{}

// Open returns a new connection to the database
func (d *KVDriver) Open(name string) (driver.Conn, error) {
	return &KVConn{}, nil
}

// KVConn implements database/sql/driver.Conn interface
type KVConn struct {
}

func (c *KVConn) Prepare(query string) (driver.Stmt, error) {
	return &KVStmt{
		conn:  c,
		query: query,
	}, nil
}

func (c *KVConn) Close() error {
	return nil
}

func (c *KVConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}

// KVStmt implements database/sql/driver.Stmt interface
type KVStmt struct {
	conn  *KVConn
	query string
}

func (s *KVStmt) Close() error {
	return nil
}

func (s *KVStmt) NumInput() int {
	return -1 // variable number of inputs
}

func (s *KVStmt) Exec(args []driver.Value) (driver.Result, error) {
	query := s.query
	trimmed := strings.TrimSpace(strings.ToUpper(query))

	if strings.HasPrefix(trimmed, "SET") {
		if len(args) != 2 {
			return nil, fmt.Errorf("SET requires exactly 2 arguments: key and value")
		}

		key, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("key must be string")
		}

		value, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("value must be string")
		}

		mu.Lock()
		kvStorage[key] = value
		mu.Unlock()

		return driver.RowsAffected(1), nil
	}

	return nil, fmt.Errorf("unsupported operation")
}

func (s *KVStmt) Query(args []driver.Value) (driver.Rows, error) {
	query := s.query
	trimmed := strings.TrimSpace(strings.ToUpper(query))

	if strings.HasPrefix(trimmed, "GET") {
		if len(args) != 1 {
			return nil, fmt.Errorf("GET requires exactly 1 argument: key")
		}

		key, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("key must be string")
		}

		mu.RLock()
		value, exists := kvStorage[key]
		mu.RUnlock()

		if !exists {
			return nil, sql.ErrNoRows
		}

		return &KVRows{value: value}, nil
	}

	return nil, fmt.Errorf("unsupported operation")
}

// KVRows implements database/sql/driver.Rows interface
type KVRows struct {
	value string
	done  bool
}

func (r *KVRows) Columns() []string {
	return []string{"value"}
}

func (r *KVRows) Close() error {
	return nil
}

func (r *KVRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	dest[0] = r.value
	return nil
}
