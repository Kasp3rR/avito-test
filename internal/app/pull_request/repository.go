package pullrequest

import (
	"avito-tech/internal/apperrors"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB interface {
	GetPool(_ context.Context) *pgxpool.Pool
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type PullRequestRepo struct {
	db DB
}

func NewRepo(db DB) *PullRequestRepo {
	return &PullRequestRepo{
		db: db,
	}
}

func (request *PullRequestRepo) create(ctx context.Context, pr *PullRequestEntity) (*PullRequestEntity, error) {
	tx, err := request.db.GetPool(ctx).Begin(ctx)
	if err != nil {
		log.Printf("[PullRequestRepo.create] failed to begin transaction: %v", err)
		return nil, apperrors.ErrDB
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var teamID uint64
	err = tx.QueryRow(ctx, `
		SELECT team_id
		FROM users
		WHERE user_id = $1
	`, pr.AuthorID).Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[PullRequestRepo.create] author's team not found for user '%s'", pr.AuthorID)
			return nil, apperrors.ErrNotFound
		}
		log.Printf("[PullRequestRepo.create] db error fetching author's team for user '%s': %v", pr.AuthorID, err)
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO pull_request (
			pull_request_id, pull_request_name, author_id, status
		) VALUES ($1, $2, $3, 'OPEN')
	`, pr.PullRequestID, pr.PullRequestName, pr.AuthorID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("[PullRequestRepo.create] PR already exists: '%s'", pr.PullRequestID)
			return nil, apperrors.ErrPRExists
		}
		log.Printf("[PullRequestRepo.create] db error inserting PR '%s': %v", pr.PullRequestID, err)
		return nil, apperrors.ErrDB
	}

	rows, err := tx.Query(ctx, `
		SELECT user_id
		FROM users
		WHERE team_id = $1
		  AND user_id <> $2
		  AND is_active = true
		LIMIT 2
	`, teamID, pr.AuthorID)
	if err != nil {
		log.Printf("[PullRequestRepo.create] db error fetching reviewers for PR '%s': %v", pr.PullRequestID, err)
		return nil, apperrors.ErrDB
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var uid string
		if err = rows.Scan(&uid); err != nil {
			log.Printf("[PullRequestRepo.create] failed to scan reviewer for PR '%s': %v", pr.PullRequestID, err)
			return nil, apperrors.ErrDB
		}
		reviewers = append(reviewers, uid)
	}

	for _, reviewerID := range reviewers {
		_, err = tx.Exec(ctx, `
			INSERT INTO pull_request_reviewer (
				pull_request_id, user_id
			) VALUES ($1, $2)
		`, pr.PullRequestID, reviewerID)
		if err != nil {
			log.Printf("[PullRequestRepo.create] failed to insert reviewer '%s' for PR '%s': %v", reviewerID, pr.PullRequestID, err)
			return nil, apperrors.ErrDB
		}
	}

	pr.Status = "OPEN"
	pr.AssignedReviewers = reviewers
	log.Printf("[PullRequestRepo.create] PR '%s' created successfully with reviewers: %v", pr.PullRequestID, reviewers)
	return pr, nil
}

func (request *PullRequestRepo) reassignReviewer(ctx context.Context, prID string, oldUserID string) (*PullRequestEntity, string, error) {
	tx, err := request.db.GetPool(ctx).Begin(ctx)
	if err != nil {
		log.Printf("[PullRequestRepo.reassignReviewer] failed to begin transaction: %v", err)
		return nil, "", apperrors.ErrDB
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var (
		status   string
		authorID string
		teamID   int64
	)

	err = tx.QueryRow(ctx, `
        SELECT pr.status, pr.author_id, u.team_id
        FROM pull_request pr
        JOIN users u ON u.user_id = pr.author_id
        WHERE pr.pull_request_id = $1
    `, prID).Scan(&status, &authorID, &teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[PullRequestRepo.reassignReviewer] PR not found: '%s'", prID)
			return nil, "", apperrors.ErrNotFound
		}
		log.Printf("[PullRequestRepo.reassignReviewer] db error fetching PR '%s': %v", prID, err)
		return nil, "", apperrors.ErrDB
	}

	if status == "MERGED" {
		log.Printf("[PullRequestRepo.reassignReviewer] cannot reassign reviewers for merged PR '%s'", prID)
		return nil, "", apperrors.ErrPRMerged
	}

	var assigned bool
	err = tx.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM pull_request_reviewer
            WHERE pull_request_id = $1 AND user_id = $2
        )
    `, prID, oldUserID).Scan(&assigned)
	if err != nil {
		log.Printf("[PullRequestRepo.reassignReviewer] db error checking assignment of user '%s' for PR '%s': %v", oldUserID, prID, err)
		return nil, "", apperrors.ErrDB
	}

	if !assigned {
		log.Printf("[PullRequestRepo.reassignReviewer] user '%s' not assigned to PR '%s'", oldUserID, prID)
		return nil, "", apperrors.ErrNotAssigned
	}

	var newUserID string
	err = tx.QueryRow(ctx, `
        SELECT user_id
        FROM users
        WHERE team_id = $1
          AND is_active = true
          AND user_id <> $2
          AND user_id <> $3
          AND NOT EXISTS (
              SELECT 1 FROM pull_request_reviewer
              WHERE pull_request_id = $4 AND user_id = users.user_id
          )
        LIMIT 1
    `, teamID, authorID, oldUserID, prID).Scan(&newUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[PullRequestRepo.reassignReviewer] no candidate available to replace user '%s' in PR '%s'", oldUserID, prID)
			return nil, "", apperrors.ErrNoCandidate
		}
		log.Printf("[PullRequestRepo.reassignReviewer] db error finding replacement for PR '%s': %v", prID, err)
		return nil, "", apperrors.ErrDB
	}

	_, err = tx.Exec(ctx, `
        DELETE FROM pull_request_reviewer
        WHERE pull_request_id = $1 AND user_id = $2
    `, prID, oldUserID)
	if err != nil {
		log.Printf("[PullRequestRepo.reassignReviewer] failed to remove old reviewer '%s' from PR '%s': %v", oldUserID, prID, err)
		return nil, "", apperrors.ErrDB
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO pull_request_reviewer(pull_request_id, user_id)
        VALUES ($1, $2)
    `, prID, newUserID)
	if err != nil {
		log.Printf("[PullRequestRepo.reassignReviewer] failed to insert new reviewer '%s' for PR '%s': %v", newUserID, prID, err)
		return nil, "", apperrors.ErrDB
	}

	pr, err := request.getByIDTx(ctx, tx, prID)
	if err != nil {
		log.Printf("[PullRequestRepo.reassignReviewer] failed to reload PR '%s' after reassignment: %v", prID, err)
		return nil, "", apperrors.ErrDB
	}

	log.Printf("[PullRequestRepo.reassignReviewer] user '%s' replaced by '%s' in PR '%s'", oldUserID, newUserID, prID)
	return pr, newUserID, nil
}

func (request *PullRequestRepo) getByIDTx(ctx context.Context, tx pgx.Tx, prID string) (*PullRequestEntity, error) {
	var pr PullRequestEntity

	err := tx.QueryRow(ctx, `
        SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
        FROM pull_request
        WHERE pull_request_id = $1
    `, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		log.Printf("[PullRequestRepo.getByIDTx] failed to fetch PR '%s': %v", prID, err)
		return nil, apperrors.ErrDB
	}

	rows, err := tx.Query(ctx, `
        SELECT user_id
        FROM pull_request_reviewer
        WHERE pull_request_id = $1
    `, prID)
	if err != nil {
		log.Printf("[PullRequestRepo.getByIDTx] failed to fetch reviewers for PR '%s': %v", prID, err)
		return nil, apperrors.ErrDB
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			log.Printf("[PullRequestRepo.getByIDTx] failed to scan reviewer for PR '%s': %v", prID, err)
			return nil, apperrors.ErrDB
		}
		reviewers = append(reviewers, uid)
	}
	pr.AssignedReviewers = reviewers

	log.Printf("[PullRequestRepo.getByIDTx] fetched PR '%s' with reviewers: %v", prID, reviewers)
	return &pr, nil
}

func (request *PullRequestRepo) merge(ctx context.Context, prID string) (*PullRequestEntity, error) {
	var entity PullRequestEntity
	// Пытаемся обновить PR если он ещё не MERGED
	err := request.db.ExecQueryRow(ctx, `
		UPDATE pull_request
		SET status = 'MERGED', merged_at = NOW()
		WHERE pull_request_id = $1 AND status <> 'MERGED'
		RETURNING 
			pull_request_id, 
			pull_request_name, 
			author_id, 
			status, 
			created_at, 
			merged_at
	`, prID).Scan(
		&entity.PullRequestID,
		&entity.PullRequestName,
		&entity.AuthorID,
		&entity.Status,
		&entity.CreatedAt,
		&entity.MergedAt,
	)

	if err == nil {
		// Успешно обновили OPEN -> MERGED
		log.Printf("[PullRequestRepo.merge] PR '%s' merged successfully", prID)
		return &entity, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		// Реальная ошибка БД
		log.Printf("[PullRequestRepo.merge] db error updating PR '%s': %v", prID, err)
		return nil, apperrors.ErrDB
	}

	// Если ничего не обновилось — значит либо PR не найден, либо уже MERGED
	err = request.db.ExecQueryRow(ctx, `
		SELECT 
			pull_request_id, 
			pull_request_name, 
			author_id, 
			status, 
			created_at, 
			merged_at
		FROM pull_request
		WHERE pull_request_id = $1
	`, prID).Scan(
		&entity.PullRequestID,
		&entity.PullRequestName,
		&entity.AuthorID,
		&entity.Status,
		&entity.CreatedAt,
		&entity.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[PullRequestRepo.merge] PR '%s' not found", prID)
			return nil, apperrors.ErrNotFound
		}

		log.Printf("[PullRequestRepo.merge] db error selecting PR '%s': %v", prID, err)
		return nil, apperrors.ErrDB
	}

	// Возвращаем текущее состояние (идемпотентность)
	log.Printf("[PullRequestRepo.merge] PR '%s' already merged, returning existing state", prID)
	return &entity, nil
}
