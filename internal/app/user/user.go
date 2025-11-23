package user

import (
	pullrequest "avito-tech/internal/app/pull_request"
	"context"
)

type Repo interface {
	getByID(ctx context.Context, id string) (*UserEntity, error)
	setIsActive(ctx context.Context, userID string, isActive bool) (*UserEntity, error)
	create(ctx context.Context, entities []*UserEntity) error
	getByTeamID(ctx context.Context, id uint64) ([]*UserEntity, error)
	getReview(ctx context.Context, userID string) ([]pullrequest.PullRequestShortDTO, error)
}

type User struct {
	repo Repo
}

func NewUser(repo Repo) *User {
	return &User{
		repo: repo,
	}
}

func (u *User) GetByID(ctx context.Context, id string) (*UserDTO, error) {
	var dto UserDTO
	entity, err := u.repo.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	dto.MapFromModel(entity)
	return &dto, err
}

func (u *User) GetByTeamID(ctx context.Context, id uint64) ([]*UserDTO, error) {
	entities, err := u.repo.getByTeamID(ctx, id)
	if err != nil {
		return nil, err
	}
	dto := MapFromModels(entities)
	return dto, nil
}

func (u *User) SetIsActive(ctx context.Context, id string, isActive bool) (*UserDTO, error) {
	var dto UserDTO
	entity, err := u.repo.setIsActive(ctx, id, isActive)
	if err != nil {
		return nil, err
	}
	dto.MapFromModel(entity)
	return &dto, err
}

func (u *User) Create(ctx context.Context, users []*UserDTO) error {
	entities := make([]*UserEntity, len(users))
	for i, v := range users {
		en := v.MapToModel()
		entities[i] = en
	}
	err := u.repo.create(ctx, entities)
	return err
}

func (u *User) GetReview(ctx context.Context, userID string) ([]*pullrequest.PullRequestShortDTOFromHttp, error) {
	entities, err := u.repo.getReview(ctx, userID)
	if err != nil {
		return nil, err
	}
	answer := pullrequest.MapToModelsShort(entities)
	return answer, nil
}
