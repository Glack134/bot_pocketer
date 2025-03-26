package domain

import "context"

type Repository interface {
	CreateUser(ctx context.Context, userID int64) error
	UserExists(ctx context.Context, userID int64) (bool, error)
}
