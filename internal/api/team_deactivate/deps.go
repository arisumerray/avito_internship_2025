package team_deactivate

import "context"

type Repository interface {
	DeactivateTeam(ctx context.Context, teamID int32) error
}
