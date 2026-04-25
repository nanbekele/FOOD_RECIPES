package db

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies the SQL migrations from the local migrations folder to the target Postgres DSN.
// DSN example: postgres://postgres:mysecretpassword@localhost:5433/foodrecipes?sslmode=disable
func RunMigrations(dsn string) error {
    // Build a robust file:// URL to the migrations directory (works on Windows/Linux/Docker)
    src := migrationsURL()
    m, err := migrate.New(src, dsn)
    if err != nil {
        return fmt.Errorf("migrate init failed: %w", err)
    }
    // Run migrations and evaluate result outside short var scope
    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migrate up failed: %w", err)
    }
    if err == migrate.ErrNoChange {
        log.Println("No database changes to apply")
    } else {
        log.Println("Database migrations applied successfully")
    }
    return nil
}

// migrationsURL returns a file:// URL for the migrations directory.
// Prefers an absolute path (to avoid CWD issues), falls back to relative.
func migrationsURL() string {
    if wd, err := os.Getwd(); err == nil {
        abs := filepath.Join(wd, "migrations")
        if _, statErr := os.Stat(abs); statErr == nil {
            return "file://" + filepath.ToSlash(abs)
        }
    }
    return "file://migrations"
}
