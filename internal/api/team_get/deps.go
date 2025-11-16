package team_get

import (
	"context"

	"avito-test-2025/internal/repository/postgresql"
)

type Repository interface {
	GetTeam(ctx context.Context, name string) (*postgresql.Team, error)
}
