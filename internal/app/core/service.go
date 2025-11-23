package core

import (
	pullrequest "avito-tech/internal/app/pull_request"
	"avito-tech/internal/app/team"
	"avito-tech/internal/app/user"
	"context"
)

type Team interface {
	GetByTeamName(ctx context.Context, teamName string) (*team.TeamDTO, error)
	Create(ctx context.Context, dto *team.TeamDTO) error
}

type User interface {
	GetByID(ctx context.Context, id string) (*user.UserDTO, error)
	GetByTeamID(ctx context.Context, id uint64) ([]*user.UserDTO, error)
	SetIsActive(ctx context.Context, id string, isActive bool) (*user.UserDTO, error)
	Create(ctx context.Context, users []*user.UserDTO) error
	GetReview(ctx context.Context, userID string) ([]*pullrequest.PullRequestShortDTOFromHttp, error)
}

type PullRequest interface {
	Create(ctx context.Context, prShort *pullrequest.PullRequestShortDTOFromHttp) (*pullrequest.PullRequestDTOFromHttp, error)
	Merge(ctx context.Context, prID string) (*pullrequest.PullRequestDTOFromHttp, error)
	Reassign(ctx context.Context, prID string, oldUserID string) (*pullrequest.PullRequestDTOFromHttp, string, error)
}

type Service struct {
	team        Team
	user        User
	pullRequest PullRequest
}

func NewService(team Team, user User, pullRequest PullRequest) *Service {
	return &Service{
		team:        team,
		user:        user,
		pullRequest: pullRequest,
	}
}
