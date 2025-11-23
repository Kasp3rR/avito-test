package team

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

type TeamRepo struct {
	db DB
}

func NewTeamRepo(db DB) *TeamRepo {
	return &TeamRepo{
		db: db,
	}
}

func (t *TeamRepo) create(ctx context.Context, teamName string, members []*TeamMemberEntity) error {
	tx, err := t.db.GetPool(ctx).BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("[TeamRepo.create] failed to begin transaction: %v", err)
		return apperrors.ErrDB
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()


	_, err = tx.Exec(ctx, `
		INSERT INTO team (team_name) VALUES ($1)
	`, teamName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("[TeamRepo.create] team with this name already exists: %v", err)
			return apperrors.ErrTeamExists
		}
		log.Printf("[TeamRepo.create] error inserting team into DB: %v", err)
		return apperrors.ErrDB
	}


	for _, member := range members {
		_, err := tx.Exec(ctx, `
			INSERT INTO users (user_id, username, team_id, is_active)
			VALUES ($1, $2, (SELECT id FROM team WHERE team_name=$3), $4)
			ON CONFLICT (user_id) DO UPDATE
			SET username=EXCLUDED.username, team_id=EXCLUDED.team_id, is_active=EXCLUDED.is_active
		`, member.UserID, member.Username, teamName, member.IsActive)
		if err != nil {
			log.Printf("[TeamRepo.create] failed to insert/update user %s: %v", member.UserID, err)
			return apperrors.ErrDB
		}
	}

	log.Printf("[TeamRepo.create] team '%s' created/updated successfully with %d members", teamName, len(members))
	return nil
}

func (t *TeamRepo) getByName(ctx context.Context, teamName string) (*TeamEntity, []TeamMemberEntity, error) {
	var entity TeamEntity
	err := t.db.Get(ctx, &entity, "SELECT id, team_name FROM team WHERE team_name=$1", teamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[TeamRepo.getByName] team not found: '%s', %v", teamName, err)
			return nil, nil, apperrors.ErrNotFound
		}
		log.Printf("[TeamRepo.getByName] DB error while fetching team '%s': %v", teamName, err)
		return nil, nil, apperrors.ErrDB
	}

	rows, err := t.db.GetPool(ctx).Query(ctx, `
        SELECT user_id, username, is_active
        FROM users
        WHERE team_id = $1
    `, entity.ID)
	if err != nil {
		log.Printf("[TeamRepo.getByName] DB error while fetching members for team '%s': %v", teamName, err)
		return nil, nil, apperrors.ErrDB
	}
	defer rows.Close()

	var members []TeamMemberEntity
	for rows.Next() {
		var m TeamMemberEntity
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			log.Printf("[TeamRepo.getByName] failed to scan member for team '%s': %v", teamName, err)
			return nil, nil, apperrors.ErrDB
		}
		members = append(members, m)
	}

	log.Printf("[TeamRepo.getByName] fetched team '%s' with %d members", teamName, len(members))
	return &entity, members, nil
}
