package postgres

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestApplyMigrationsIsSafeForConcurrentCallers is a regression test for a bug
// where ApplyMigrations ran "CREATE SCHEMA IF NOT EXISTS system" (and later,
// per-migration DDL) directly against a *pgxpool.Pool before acquiring its
// advisory lock. pg_advisory_lock is session-scoped, but pool.Exec/Begin each
// borrow a (possibly different) connection per call, so the lock provided no
// real mutual exclusion: two concurrent callers - e.g. two replicas of
// cmd/migrate racing on a fresh database, or this repository's own
// postgres_integration_test.go and store_integration_test.go running in
// parallel via `go test ./...` - could both attempt
// "CREATE SCHEMA IF NOT EXISTS system" at once and fail with
// "duplicate key value violates unique constraint pg_namespace_nspname_index".
//
// ApplyMigrations now acquires one dedicated connection up front and runs the
// whole procedure (lock, journal creation, migration loop, unlock) on it.
// This test drives several concurrent callers against one fresh database and
// asserts all of them succeed and agree on the resulting schema.
//
// Set TEST_DATABASE_URL to run it; it is skipped otherwise.
func TestApplyMigrationsIsSafeForConcurrentCallers(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping PostgreSQL integration test")
	}
	ctx := context.Background()

	const concurrentCallers = 8
	var wg sync.WaitGroup
	errs := make([]error, concurrentCallers)
	for i := 0; i < concurrentCallers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pool, err := pgxpool.New(ctx, url)
			if err != nil {
				errs[i] = err
				return
			}
			defer pool.Close()
			errs[i] = ApplyMigrations(ctx, pool)
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("concurrent caller %d: %v", i, err)
		}
	}

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatalf("connect to verify: %v", err)
	}
	defer pool.Close()
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM system.schema_migration`).Scan(&count); err != nil {
		t.Fatalf("count applied migrations: %v", err)
	}
	migrations, err := LoadMigrations()
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}
	if count != len(migrations) {
		t.Fatalf("expected %d recorded migrations, got %d", len(migrations), count)
	}
}
