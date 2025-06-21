package stdlib

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
)

// This stub registers a minimal PostgreSQL driver named "pgx" to allow
// building without the real github.com/jackc/pgx/v5/stdlib package.
func init() {
	sql.Register("pgx", &stubDriver{})
}

type stubDriver struct{}

func (d *stubDriver) Open(name string) (driver.Conn, error) {
	return &stubConn{}, nil
}

type stubConn struct{}

func (c *stubConn) Prepare(query string) (driver.Stmt, error) { return &stubStmt{}, nil }
func (c *stubConn) Close() error                              { return nil }
func (c *stubConn) Begin() (driver.Tx, error)                 { return &stubTx{}, nil }

// Implement driver.Pinger so db.Ping() succeeds.
func (c *stubConn) Ping(ctx context.Context) error { return nil }

// Optional interfaces with no-op implementations.
func (c *stubConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return stubResult{}, nil
}
func (c *stubConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return &stubRows{}, nil
}

// Statement implementation.
type stubStmt struct{}

func (s *stubStmt) Close() error                                    { return nil }
func (s *stubStmt) NumInput() int                                   { return -1 }
func (s *stubStmt) Exec(args []driver.Value) (driver.Result, error) { return stubResult{}, nil }
func (s *stubStmt) Query(args []driver.Value) (driver.Rows, error)  { return &stubRows{}, nil }

// Transaction implementation.
type stubTx struct{}

func (tx *stubTx) Commit() error   { return nil }
func (tx *stubTx) Rollback() error { return nil }

// Result implementation.
type stubResult struct{}

func (r stubResult) LastInsertId() (int64, error) { return 0, nil }
func (r stubResult) RowsAffected() (int64, error) { return 0, nil }

// Rows implementation.
type stubRows struct{}

func (r *stubRows) Columns() []string              { return []string{} }
func (r *stubRows) Close() error                   { return nil }
func (r *stubRows) Next(dest []driver.Value) error { return io.EOF }
