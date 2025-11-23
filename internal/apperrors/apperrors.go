package apperrors

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrDB          = errors.New("error with db")
	ErrTeamExists  = errors.New("team already exists")
	ErrPRExists    = errors.New("pr already exists")
	ErrPRMerged    = errors.New("pr already merged")
	ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
)
