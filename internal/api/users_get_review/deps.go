package users_get_review

import (
	"context"
	
	"avito-test-2025/internal/repository/postgresql"
)

type Repository interface {
	GetByReviewerID(ctx context.Context, reviewerID string) ([]postgresql.PR, error)
}
