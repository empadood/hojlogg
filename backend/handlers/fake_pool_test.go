package handlers_test

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// fakePool satisfies handlers.Pooler but returns errors / empty results.
type fakePool struct{}

func (f *fakePool) Ping(_ context.Context) error { return nil }

func (f *fakePool) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	// RowsAffected() == 0 → handlers return 404
	return pgconn.NewCommandTag("DELETE 0"), nil
}

func (f *fakePool) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return &fakeRows{}, nil
}

func (f *fakePool) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return &fakeRow{}
}

// ── fakeRows ─────────────────────────────────────────────────────────────────

type fakeRows struct{ closed bool }

func (r *fakeRows) Next() bool                                      { return false }
func (r *fakeRows) Scan(_ ...any) error                             { return nil }
func (r *fakeRows) Values() ([]any, error)                          { return nil, nil }
func (r *fakeRows) RawValues() [][]byte                             { return nil }
func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription    { return nil }
func (r *fakeRows) Close()                                          { r.closed = true }
func (r *fakeRows) Err() error                                      { return nil }
func (r *fakeRows) Conn() *pgx.Conn                                 { return nil }
func (r *fakeRows) CommandTag() pgconn.CommandTag                   { return pgconn.CommandTag{} }

// ── fakeRow ───────────────────────────────────────────────────────────────────

type fakeRow struct{}

func (r *fakeRow) Scan(_ ...any) error {
	return fmt.Errorf("no rows in result set")
}
