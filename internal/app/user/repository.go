package user

import (
	pullrequest "avito-tech/internal/app/pull_request"
	"avito-tech/internal/apperrors"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type DB interface {
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type UserRepo struct {
	db DB
}

func NewUserRepo(db DB) *UserRepo {
	return &UserRepo{db: db}
}

func (user *UserRepo) getByID(ctx context.Context, id string) (*UserEntity, error) {
	var entity UserEntity
	err := user.db.Get(ctx, &entity, "SELECT user_id, username, team_id, is_active FROM users WHERE user_id=$1", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[UserRepo.getByID] user '%s' not found", id)
			return nil, apperrors.ErrNotFound
		}
		log.Printf("[UserRepo.getByID] db error fetching user '%s': %v", id, err)
		return nil, apperrors.ErrDB
	}
	log.Printf("[UserRepo.getByID] fetched user '%s'", id)
	return &entity, nil
}

func (user *UserRepo) getByTeamID(ctx context.Context, id uint64) ([]*UserEntity, error) {
	var entities []*UserEntity
	err := user.db.Select(ctx, &entities, "SELECT user_id, username, team_id, is_active FROM users WHERE team_id=$1", id)
	if err != nil {
		log.Printf("[UserRepo.getByTeamID] db error fetching users for team '%d': %v", id, err)
		return nil, apperrors.ErrDB
	}
	if len(entities) == 0 {
		log.Printf("[UserRepo.getByTeamID] no users found for team '%d'", id)
		return nil, apperrors.ErrNotFound
	}
	log.Printf("[UserRepo.getByTeamID] fetched %d users for team '%d'", len(entities), id)
	return entities, nil
}

func (user *UserRepo) setIsActive(ctx context.Context, userID string, isActive bool) (*UserEntity, error) {
	var entity UserEntity

	err := user.db.ExecQueryRow(ctx, `
		UPDATE users u
		SET is_active = $1
		FROM team t
		WHERE u.user_id = $2
		  AND t.id = u.team_id
		RETURNING 
			u.user_id, 
			u.username, 
			u.team_id, 
			t.team_name,
			u.is_active
	`, isActive, userID).Scan(
		&entity.UserID,
		&entity.Username,
		&entity.TeamID,
		&entity.TeamName,
		&entity.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[UserRepo.setIsActive] user '%s' not found", userID)
			return nil, apperrors.ErrNotFound
		}
		log.Printf("[UserRepo.setIsActive] db error updating user '%s': %v", userID, err)
		return nil, apperrors.ErrDB
	}

	log.Printf("[UserRepo.setIsActive] updated user '%s' isActive=%v", userID, isActive)
	return &entity, nil
}

func (user *UserRepo) create(ctx context.Context, entities []*UserEntity) error {
	values := []interface{}{}
	placeholders := []string{}

	for i, u := range entities {
		n := i * 4
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", n+1, n+2, n+3, n+4))
		values = append(values, u.UserID, u.Username, u.TeamID, u.IsActive)
	}

	query := fmt.Sprintf(`
        INSERT INTO users (user_id, username, team_id, is_active)
        VALUES %s
        ON CONFLICT (user_id) DO UPDATE
        SET username = EXCLUDED.username,
            team_id = EXCLUDED.team_id,
            is_active = EXCLUDED.is_active
    `, strings.Join(placeholders, ","))

	_, err := user.db.Exec(ctx, query, values...)
	if err != nil {
		log.Printf("[UserRepo.create] db error inserting/updating users: %v", err)
		return apperrors.ErrDB
	}

	log.Printf("[UserRepo.create] inserted/updated %d users", len(entities))
	return nil
}

func (user *UserRepo) getReview(ctx context.Context, userID string) ([]pullrequest.PullRequestShortDTO, error) {
	var exists bool
	err := user.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)", userID)
	if err != nil {
		log.Printf("[UserRepo.getReview] db error checking existence of user '%s': %v", userID, err)
		return nil, apperrors.ErrDB
	}
	if !exists {
		log.Printf("[UserRepo.getReview] user '%s' not found", userID)
		return nil, apperrors.ErrNotFound
	}

	var prs []pullrequest.PullRequestShortDTO
	query := `
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.status
		FROM pull_request pr
		JOIN pull_request_reviewer prr
			ON pr.pull_request_id = prr.pull_request_id
		WHERE prr.user_id = $1
	`
	err = user.db.Select(ctx, &prs, query, userID)
	if err != nil {
		log.Printf("[UserRepo.getReview] db error fetching PRs for reviewer '%s': %v", userID, err)
		return nil, apperrors.ErrDB
	}

	log.Printf("[UserRepo.getReview] fetched %d PRs for reviewer '%s'", len(prs), userID)
	return prs, nil
}
