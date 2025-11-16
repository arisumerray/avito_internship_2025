package api

import (
	"context"

	"avito-test-2025/internal/repository/postgresql"
)

type Repository interface {
	Reassign(ctx context.Context, prID string, oldUserID string) (pr *postgresql.PR, newUserID string, err error)
	AddTeam(ctx context.Context, name string, members []postgresql.User) error
	GetTeam(ctx context.Context, name string) (*postgresql.Team, error)
	Merge(ctx context.Context, prID string) (*postgresql.PR, error)
	Create(ctx context.Context, prID, prName, authorID string) (*postgresql.PR, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) (*postgresql.User, error)
	GetByReviewerID(ctx context.Context, reviewerID string) ([]postgresql.PR, error)
	DeactivateTeam(ctx context.Context, teamID int32) error
}
