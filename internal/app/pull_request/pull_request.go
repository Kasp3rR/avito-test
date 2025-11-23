package pullrequest

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Repo interface {
	create(ctx context.Context, pr *PullRequestEntity) (*PullRequestEntity, error)
	reassignReviewer(ctx context.Context, prID string, oldUserID string) (*PullRequestEntity, string, error)
	getByIDTx(ctx context.Context, tx pgx.Tx, prID string) (*PullRequestEntity, error)
	merge(ctx context.Context, prID string) (*PullRequestEntity, error)
}

type PullRequest struct {
	repo Repo
}

func NewPullRequest(repo Repo) *PullRequest {
	return &PullRequest{
		repo: repo,
	}
}

func (pr *PullRequest) Create(ctx context.Context, prShort *PullRequestShortDTOFromHttp) (*PullRequestDTOFromHttp, error) {
	entity, err := pr.repo.create(ctx, prShort.MapToPREntity())
	if err != nil {
		return nil, err
	}
	var answer PullRequestDTOFromHttp
	answer.MapFromModel(entity)
	return &answer, nil
}

func (pr *PullRequest) Merge(ctx context.Context, prID string) (*PullRequestDTOFromHttp, error) {
	entity, err := pr.repo.merge(ctx, prID)
	if err != nil {
		return nil, err
	}
	var answer PullRequestDTOFromHttp
	answer.MapFromModel(entity)
	return &answer, nil
}

func (pr *PullRequest) Reassign(ctx context.Context, prID string, oldUserID string) (*PullRequestDTOFromHttp, string, error) {
	entity, newUser, err := pr.repo.reassignReviewer(ctx, prID, oldUserID)
	if err != nil {
		return nil, "", err
	}
	var answer PullRequestDTOFromHttp
	answer.MapFromModel(entity)
	return &answer, newUser, err
}
