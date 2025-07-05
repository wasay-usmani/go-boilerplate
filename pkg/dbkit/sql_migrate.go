package dbkit

import (
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/wasay-usmani/go-boilerplate/pkg/logkit"
)

const (
	_migrations = "migrations"
	_schema     = "schema"
	_app        = "app"
)

func ApplyMigrations(superUserDBURL, appDBURL, dbSchema string, l logkit.Logger, schemaAssets, appAssets embed.FS) error {
	// Init schema db connection
	schemaDB, err := LoadMySqlConn(superUserDBURL, true, l)
	if err != nil {
		return fmt.Errorf("failed to load superuser db, error: %w", err)
	}
	defer schemaDB.Close()

	// Apply schema migrations
	if err = applySchemaMigrations(schemaDB.DB, dbSchema, l, schemaAssets); err != nil {
		return fmt.Errorf("failed to apply schema migrations, error: %w", err)
	}

	// Init app db connection
	appDB, err := LoadMySqlConn(appDBURL, true, l)
	if err != nil {
		return fmt.Errorf("failed to load sfpy_reporter, error: %w", err)
	}
	defer appDB.Close()

	// Apply app migrations
	if err = applyAppMigrations(appDB.DB, dbSchema, l, appAssets); err != nil {
		return fmt.Errorf("failed to apply app migrations, error: %w", err)
	}

	return nil
}

func applySchemaMigrations(db *sql.DB, dbSchema string, l logkit.Logger, schemaAssets embed.FS) error {
	migrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: schemaAssets,
		Root:       "migrations",
	}

	migrationSet := &migrate.MigrationSet{
		TableName: fmt.Sprintf("%s_%s_%s", _migrations, _schema, dbSchema),
	}

	n, err := migrationSet.Exec(db, "mysql", migrationSource, migrate.Up)
	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("applied %d schema migrations", n))
	return nil
}

func applyAppMigrations(db *sql.DB, dbSchema string, l logkit.Logger, appAssets embed.FS) error {
	migrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: appAssets,
		Root:       "migrations",
	}

	migrationSet := &migrate.MigrationSet{
		TableName: fmt.Sprintf("%s_%s_%s", _migrations, _app, dbSchema),
	}

	n, err := migrationSet.Exec(db, "mysql", migrationSource, migrate.Up)
	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("applied %d app migrations", n))
	return nil
}
