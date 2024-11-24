package sqlx

import (
	"context"
	"database/sql"
)

func (tx *Tx) Commit() error {
	return tx.TX.Commit()
}

func (tx *Tx) Exec(query string, args ...any) (sql.Result, error) {
	return tx.TX.Exec(query, args...)
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.TX.ExecContext(ctx, query, args...)
}

func (tx *Tx) Prepare(query string) (*Stmt, error) {
	stmt, err := tx.TX.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{Stmt: stmt}, nil
}

func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := tx.TX.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{Stmt: stmt}, nil
}

func (tx *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	return tx.TX.Query(query, args...)
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.TX.QueryContext(ctx, query, args...)
}

func (tx *Tx) QueryRow(query string, args ...any) *sql.Row {
	return tx.TX.QueryRow(query, args...)
}

func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return tx.TX.QueryRowContext(ctx, query, args...)
}

func (tx *Tx) Rollback() error {
	return tx.TX.Rollback()
}

func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	return &Stmt{Stmt: tx.TX.Stmt(stmt.Stmt)}
}

func (tx *Tx) StmtContext(ctx context.Context, stmt *Stmt) *Stmt {
	return &Stmt{Stmt: tx.TX.StmtContext(ctx, stmt.Stmt)}
}
