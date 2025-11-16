package team_add

import (
	"avito-test-2025/internal/repository/postgresql"
	"context"
)

type Repository interface {
	AddTeam(ctx context.Context, name string, members []postgresql.User) error
}
