package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

func (r *Repository) Reassign(ctx context.Context, prID string, oldUserID string) (pr *PR, newUserID string, err error) {
	var teamID int32
	err = r.DB.QueryRowContext(ctx, "SELECT team_id FROM users WHERE id = $1", oldUserID).Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			// 404 PR или user не найдены
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}
	res := PR{}
	err = r.DB.QueryRowContext(ctx, "SELECT name, author_id, status, reviewer1_id, reviewer2_id FROM pr WHERE id = $1", prID).
		Scan(&res.Name, &res.AuthorID, &res.Status, &res.Reviewer1ID, &res.Reviewer2ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			// 404 PR или user не найдены
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}
	if res.Status == "MERGED" {
		return nil, "", ErrPRIsMerged
	}
	if oldUserID != res.Reviewer1ID && oldUserID != res.Reviewer2ID {
		return nil, "", ErrNotAssigned
	}
	rows, err := r.DB.QueryContext(ctx, "SELECT id, name, is_active FROM users WHERE team_id = $1 AND is_active = true", teamID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query: %w", err)
	}

	var candidates []User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Name, &u.IsActive)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan: %w", err)
		}
		if u.ID != res.Reviewer1ID && u.ID != res.Reviewer2ID && u.ID != res.AuthorID {
			candidates = append(candidates, u)
		}
	}
	if len(candidates) == 0 {
		return nil, "", ErrNoCandidates
	}
	for _, u := range candidates {
		tx, err := r.DB.BeginTx(ctx, nil)
		if err != nil {
			return nil, "", err
		}
		err = tx.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1 FOR UPDATE SKIP LOCKED", u.ID).Scan(&u.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
				continue
			}
			tx.Rollback()
			return nil, "", fmt.Errorf("failed to query: %w", err)
		}
		if res.Reviewer1ID == u.ID || res.Reviewer2ID == u.ID || res.AuthorID == u.ID {
			continue
		}
		err = tx.QueryRowContext(ctx, "SELECT id, status FROM pr WHERE ID = $1 FOR UPDATE SKIP LOCKED", prID).Scan(&res.ID, &res.Status)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
				return nil, "", ErrNotAssigned
			}
			return nil, "", fmt.Errorf("failed to query: %w", err)
		}
		if res.Status == "MERGED" {
			tx.Rollback()
			return nil, "", ErrPRIsMerged
		}
		if res.Reviewer1ID == oldUserID {
			_, err = tx.ExecContext(ctx, "UPDATE pr SET reviewer1_id = $1 WHERE id = $2", u.ID, res.ID)
		} else if res.Reviewer2ID == oldUserID {
			_, err = tx.ExecContext(ctx, "UPDATE pr SET reviewer2_id = $1 WHERE id = $2", u.ID, res.ID)
		}
		if err != nil {
			tx.Rollback()
			return nil, "", fmt.Errorf("failed to update: %w", err)
		}
		err = tx.Commit()
		if err != nil {
			return nil, "", fmt.Errorf("failed to commit: %w", err)
		}
		return pr, u.ID, nil
	}
	return nil, "", ErrNotAssigned
}

func (r *Repository) Create(ctx context.Context, prID, prName, authorID string) (*PR, error) {
	var teamID int32
	err := r.DB.QueryRowContext(ctx, "SELECT team_id FROM users WHERE id = $1", authorID).Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			// Автор/команда не найдены
			return nil, ErrNotFound
		}
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, "SELECT id, name, is_active FROM users WHERE team_id = $1 AND is_active = true", teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var candidates []User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Name, &u.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		if u.ID != authorID {
			candidates = append(candidates, u)
		}
	}

	pr := &PR{
		ID:       prID,
		Name:     prName,
		AuthorID: authorID,
		Status:   "OPEN",
	}

	if len(candidates) == 0 {
		_, err = r.DB.ExecContext(ctx, `INSERT INTO pr (id, name, author_id, status, created_at) VALUES ($1, $2, $3, 'OPEN', $4)`, prID, prName, authorID, time.Now())
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" { // unique violation
					return nil, ErrPRExists
				}
			}
			return nil, err
		}
		return pr, nil
	}
	if len(candidates) == 1 {
		_, err = r.DB.ExecContext(ctx, `INSERT INTO pr (id, name, author_id, status, created_at, reviewer1_id) VALUES ($1, $2, $3, 'OPEN', $4, $5)`, prID, prName, authorID, time.Now(), candidates[0].ID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" { // unique violation
					return nil, ErrPRExists
				}
			}
			return nil, err
		}
		pr.Reviewer1ID = candidates[0].ID
		return pr, nil
	}
	_, err = r.DB.ExecContext(ctx, `INSERT INTO pr (id, name, author_id, status, created_at, reviewer1_id, reviewer2_id) VALUES ($1, $2, $3, 'OPEN', $4, $5, $6)`, prID, prName, authorID, time.Now(), candidates[0].ID, candidates[1].ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique violation
				return nil, ErrPRExists
			}
		}
		return nil, err
	}
	pr.Reviewer1ID = candidates[0].ID
	pr.Reviewer2ID = candidates[1].ID
	return pr, nil
}

func (r *Repository) Merge(ctx context.Context, prID string) (*PR, error) {
	var pr PR
	pr.ID = prID
	err := r.DB.QueryRowContext(ctx, "SELECT name, status, author_id, pr.reviewer1_id, pr.reviewer2_id, pr.merged_at FROM pr WHERE id = $1", prID).
		Scan(&pr.Status, &pr.AuthorID, &pr.Reviewer1ID, &pr.Reviewer2ID, &pr.MergedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
	}
	if pr.Status == "MERGED" {
		return &pr, nil
	}
	mergedAt := time.Now()
	_, err = r.DB.ExecContext(ctx, "UPDATE pr SET status = 'MERGED', merged_at = $1 WHERE id = $2", mergedAt, prID)
	if err != nil {
		return nil, err
	}
	pr.Status = "MERGED"
	pr.MergedAt = &mergedAt
	return &pr, nil
}

func (r *Repository) GetByReviewerID(ctx context.Context, reviewerID string) ([]PR, error) {
	var prs []PR
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, author_id, status FROM pr WHERE (reviewer2_id = $1 OR reviewer1_id = $1) AND status = 'OPEN'`, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var pr PR
		err = rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		prs = append(prs, pr)
	}
	return prs, nil
}
