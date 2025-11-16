package pull_request_reassign

import (
	"context"

	repo "avito-test-2025/internal/repository/postgresql"
)

type Repository interface {
	Reassign(ctx context.Context, prID string, oldUserID string) (pr *repo.PR, newUserID string, err error)
}
