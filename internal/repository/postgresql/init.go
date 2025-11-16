package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewRepository(cfg Config) (*Repository, error) {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	r := &Repository{
		DB: db,
	}

	return r, nil
}
