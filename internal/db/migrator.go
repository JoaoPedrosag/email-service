package db

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB, dir string) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			executed_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`); err != nil {
		return err
	}

	entries := make([]string, 0)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    if !d.IsDir() && filepath.Ext(path) == ".sql" && strings.HasSuffix(d.Name(), ".up.sql") {
        entries = append(entries, path)
    }
    return nil
	})

	if err != nil {
		return err
	}

	sort.Strings(entries)

	for _, file := range entries {
		version := filepath.Base(file)

		var exists string
		query := `SELECT version FROM schema_migrations WHERE version = $1`
		err := db.Get(&exists, query, version)
		if err == nil {
			continue
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec(
			`INSERT INTO schema_migrations (version) VALUES ($1)`,
			version,
		); err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		log.Printf("Applied migration %s", version)
	}

	return nil
}
