package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Ping(ctx context.Context) error {
	return errors.Wrap(s.db.Ping(ctx), "db.Ping")
}
