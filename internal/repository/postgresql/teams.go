package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) AddTeam(ctx context.Context, name string, members []User) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var teamID int64
	err = tx.QueryRowContext(ctx, `
        INSERT INTO teams (name) 
        VALUES ($1)
        RETURNING id
    `, name).Scan(&teamID)

	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				return ErrTeamExists
			}
		}
		return fmt.Errorf("failed to add team: %w", err)
	}
	query := `INSERT INTO users (id, name, is_active, team_id) VALUES `
	var args []any

	for i, u := range members {
		idx := i*4 + 1
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", idx, idx+1, idx+2, idx+3)
		if i < len(members)-1 {
			query += ", "
		}
		args = append(args, u.ID, u.Name, u.IsActive, teamID)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create users: %w, query: %s", err, query)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (r *Repository) GetTeam(ctx context.Context, name string) (*Team, error) {
	var teamID int32
	err := r.DB.QueryRowContext(ctx, "SELECT id FROM teams WHERE name = $1", name).Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) && errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	rows, err := r.DB.QueryContext(ctx, "SELECT id, name, is_active FROM users WHERE team_id = $1", teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	defer rows.Close()
	var team Team
	for rows.Next() {
		var member User
		if err := rows.Scan(&member.ID, &member.Name, &member.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}
		team.Members = append(team.Members, member)
	}
	return &team, nil
}
