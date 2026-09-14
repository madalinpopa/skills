package pgkit

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func MigrateDatabaseUp(ctx context.Context, schema string, pool *pgxpool.Pool, migrations fs.FS, dir string) error {
	db := stdlib.OpenDBFromPool(pool)
	defer func() { _ = db.Close() }()

	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %q", schema)); err != nil {
		return fmt.Errorf("creating schema %q: %w", schema, err)
	}

	source, err := iofs.New(migrations, dir)
	if err != nil {
		return fmt.Errorf("reading migrations for %q: %w", schema, err)
	}

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{
		SchemaName:      schema,
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("opening migration driver for %q: %w", schema, err)
	}

	migration, err := migrate.NewWithInstance("iofs", source, "pgx", driver)
	if err != nil {
		return fmt.Errorf("preparing migrations for %q: %w", schema, err)
	}
	defer func() {
		if sourceErr, dbErr := migration.Close(); sourceErr != nil || dbErr != nil {
			slog.Error("closing migrations", "schema", schema, "source", sourceErr, "database", dbErr)
		}
	}()

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations for %q: %w", schema, err)
	}

	return nil
}
