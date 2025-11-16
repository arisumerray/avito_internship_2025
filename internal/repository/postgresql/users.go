package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) SetIsActive(ctx context.Context, userID string, isActive bool) (*User, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	var user User
	err = tx.QueryRowContext(ctx, "SELECT name, is_active, team_id FROM users WHERE id = $1 FOR UPDATE", userID).
		Scan(&user.Name, &user.IsActive, &user.TeamID)

	if err != nil {
		tx.Rollback()
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET is_active=$1 WHERE id=$2", isActive, userID)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to set is active: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRowContext(ctx, "SELECT name FROM teams WHERE id = $1", user.TeamID).Scan(&user.TeamName)
	return &user, err
}

func (r *Repository) DeactivateTeam(ctx context.Context, teamID int32) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.QueryContext(ctx, "SELECT id FROM users WHERE team_id = $1 FOR UPDATE", teamID)

	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET is_active=false WHERE team_id=$1", teamID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set is active: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
