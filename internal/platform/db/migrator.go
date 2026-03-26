package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migrator adalah migration runner sederhana.
// Ia menjalankan file SQL manual (.up.sql / .down.sql) tanpa AutoMigrate.
type Migrator struct {
	db *sql.DB
}

// NewMigrator membuat instance migration runner.
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// Up menjalankan semua file .up.sql yang belum pernah dieksekusi.
func (m *Migrator) Up(migrationsDir string) error {
	if err := m.ensureTable(); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		applied, err := m.isApplied(filepath.Base(file))
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		tx, err := m.db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("gagal menjalankan migration %s: %w", file, err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (filename) VALUES (?)`, filepath.Base(file)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("gagal mencatat migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// DownLast membatalkan migration terakhir satu langkah.
// Ini sederhana tapi cukup untuk belajar.
func (m *Migrator) DownLast(migrationsDir string) error {
	if err := m.ensureTable(); err != nil {
		return err
	}

	var filename string
	err := m.db.QueryRow(`
		SELECT filename
		FROM schema_migrations
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	downFile := strings.Replace(filename, ".up.sql", ".down.sql", 1)
	fullPath := filepath.Join(migrationsDir, downFile)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("gagal membaca down migration %s: %w", fullPath, err)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(content)); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`DELETE FROM schema_migrations WHERE filename = ?`, filename); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (m *Migrator) ensureTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (m *Migrator) isApplied(filename string) (bool, error) {
	var count int
	err := m.db.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE filename = ?`, filename).Scan(&count)
	return count > 0, err
}
