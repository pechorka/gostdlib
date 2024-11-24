package sqlx

import (
	"context"
	"database/sql"
)

func (stmt *Stmt) Close() error {
	return stmt.Stmt.Close()
}

func (stmt *Stmt) Exec(args ...any) (sql.Result, error) {
	return stmt.Stmt.Exec(args...)
}

func (stmt *Stmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return stmt.Stmt.ExecContext(ctx, args...)
}

func (stmt *Stmt) Query(args ...any) (*sql.Rows, error) {
	return stmt.Stmt.Query(args...)
}

func (stmt *Stmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	return stmt.Stmt.QueryContext(ctx, args...)
}

func (stmt *Stmt) QueryRow(args ...any) *sql.Row {
	return stmt.Stmt.QueryRow(args...)
}

func (stmt *Stmt) QueryRowContext(ctx context.Context, args ...any) *sql.Row {
	return stmt.Stmt.QueryRowContext(ctx, args...)
}
