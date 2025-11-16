package pull_request_create

import (
	"context"

	"avito-test-2025/internal/repository/postgresql"
)

type Repository interface {
	Create(ctx context.Context, prID, prName, authorID string) (*postgresql.PR, error)
}
