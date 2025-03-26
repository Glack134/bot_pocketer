package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/polyk005/tg_bot/internal/domain"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (domain.Repository, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{db: pool}, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users(id) VALUES($1) ON CONFLICT DO NOTHING", userID)
	return err
}

func (r *PostgresRepository) UserExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	return exists, err
}
