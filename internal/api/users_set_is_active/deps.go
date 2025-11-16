package users_set_is_active

import (
	"context"

	"avito-test-2025/internal/repository/postgresql"
)

type Repository interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*postgresql.User, error)
}
