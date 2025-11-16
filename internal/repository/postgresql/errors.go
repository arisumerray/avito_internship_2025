package postgresql

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrNoCandidates = errors.New("no candidates")
	ErrNotAssigned  = errors.New("not assigned")
	ErrPRIsMerged   = errors.New("PR is merged")
	ErrPRExists     = errors.New("PR exists")
	ErrTeamExists   = errors.New("team exists")
)
