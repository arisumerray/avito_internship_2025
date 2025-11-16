package pull_request_merge

import (
	"context"
	
	"avito-test-2025/internal/repository/postgresql"
)

type Repository interface {
	Merge(ctx context.Context, prID string) (*postgresql.PR, error)
}
