package postgresql

import "time"

type User struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
	TeamID   int32  `db:"team_id"`
	TeamName string
}

type Team struct {
	ID      int32  `db:"id"`
	Name    string `db:"name"`
	Members []User
}

type PR struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	AuthorID    string     `db:"author_id"`
	Status      string     `db:"status"`
	Reviewer1ID string     `db:"reviewer1_id"`
	Reviewer2ID string     `db:"reviewer2_id"`
	CreatedAt   time.Time  `db:"created_at"`
	MergedAt    *time.Time `db:"updated_at"`
}
