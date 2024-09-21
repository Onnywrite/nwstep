package storage

import (
	"errors"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

type PgStorage struct {
	db *sqlx.DB
}

func Connect(conn string) (*PgStorage, error) {
	db, err := sqlx.Connect("pgx", conn)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New("file:///migrations", conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &PgStorage{
		db: db,
	}, nil
}

func (pg *PgStorage) Disconnect() error {
	err := pg.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
