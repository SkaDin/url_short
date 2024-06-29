package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"url_short/internal/storage"
)

type Storage struct {
	db *pgx.Conn
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgx.Connect(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	createTable := `
		CREATE TABLE IF NOT EXISTS url(
		    id SERIAL PRIMARY KEY ,
		    alias TEXT NOT NULL ,
		    url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`
	_, err = db.Exec(ctx, createTable)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveULR"
	ctx := context.Background()
	res, err := s.db.Exec(ctx, `INSERT INTO url(url, alias) VALUES ($1, $2)`, urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id := res.RowsAffected()

	if id != 1 {
		return 0, fmt.Errorf("%s: error getting number of affected rows: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	ctx := context.Background()
	var resURL string
	err := s.db.QueryRow(ctx, `SELECT url FROM url WHERE alias=$1`, alias).Scan(&resURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement %w", op, err)
	}
	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	ctx := context.Background()

	err := s.db.QueryRow(ctx, `DELETE FROM url WHERE alias=$1`, alias)
	if err != nil {
		return fmt.Errorf("%s: removed statement %v", op, err)
	}
	return nil
}
