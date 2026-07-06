package database

import (
	"fmt"
	"time"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(url string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil { return nil, fmt.Errorf("postgres: %w", err) }
	db.SetMaxOpenConns(25); db.SetMaxIdleConns(10); db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

func RunMigrations(url string) error {
	m, err := migrate.New("file://migrations", url)
	if err != nil { return err }
	if err := m.Up(); err != nil && err != migrate.ErrNoChange { return err }
	return nil
}
