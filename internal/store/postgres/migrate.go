package postgres

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Migration struct {
	Version  int
	Name     string
	SQL      string
	Checksum string
}

var canonicalMigrationNames = []string{
	"000001_extensions_and_schemas.sql",
	"000002_canonical_registry.sql",
	"000003_canonical_relationships.sql",
	"000004_isolation_profiles.sql",
	"000005_engine_topology.sql",
	"000006_markets.sql",
	"000007_digital_estates.sql",
	"000008_capabilities.sql",
	"000009_mapping_scopes.sql",
	"000010_external_references.sql",
	"000011_canonical_mappings.sql",
	"000012_capability_bindings.sql",
	"000013_context_snapshots.sql",
	"000014_audit.sql",
	"000015_messaging.sql",
	"000016_idempotency_and_revisions.sql",
	"000017_indexes_and_integrity.sql",
}

func LoadMigrations() ([]Migration, error) {
	migrations := make([]Migration, 0, len(canonicalMigrationNames))
	for _, name := range canonicalMigrationNames {
		contents, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", name, err)
		}
		version, err := strconv.Atoi(name[:6])
		if err != nil {
			return nil, fmt.Errorf("parse migration version %s: %w", name, err)
		}
		sum := sha256.Sum256(contents)
		migrations = append(migrations, Migration{Version: version, Name: name, SQL: string(contents), Checksum: hex.EncodeToString(sum[:])})
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].Version < migrations[j].Version })
	return migrations, nil
}

func ApplyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return errors.New("migration database pool is nil")
	}
	migrations, err := LoadMigrations()
	if err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS system; CREATE TABLE IF NOT EXISTS system.schema_migration (version integer PRIMARY KEY, name text NOT NULL, checksum text NOT NULL, applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		return fmt.Errorf("create migration journal: %w", err)
	}
	if _, err := pool.Exec(ctx, `SELECT pg_advisory_lock(894217301)`); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	defer pool.Exec(context.Background(), `SELECT pg_advisory_unlock(894217301)`)
	for _, migration := range migrations {
		var recordedName, checksum string
		err := pool.QueryRow(ctx, `SELECT name, checksum FROM system.schema_migration WHERE version = $1`, migration.Version).Scan(&recordedName, &checksum)
		if err == nil {
			if recordedName != migration.Name || checksum != migration.Checksum {
				return fmt.Errorf("migration %s checksum mismatch", migration.Name)
			}
			continue
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("read migration journal for %s: %w", migration.Name, err)
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", migration.Name, err)
		}
		if _, err = tx.Exec(ctx, migration.SQL); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("apply migration %s: %w", migration.Name, err)
		}
		if _, err = tx.Exec(ctx, `INSERT INTO system.schema_migration(version, name, checksum) VALUES ($1, $2, $3)`, migration.Version, migration.Name, migration.Checksum); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record migration %s: %w", migration.Name, err)
		}
		if err = tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", migration.Name, err)
		}
	}
	return nil
}

func (m Migration) String() string {
	return strings.TrimSpace(fmt.Sprintf("%06d %s", m.Version, m.Name))
}
