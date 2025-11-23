package routing

import (
	"avito-tech/internal/app/core"
	"context"
)

type ImplInterface interface {
	AddTeam(ctx context.Context, req *core.AddTeamRequest) (*core.AddTeamResponse, error)
	CreatePullRequestFromCreateRequest(ctx context.Context, request *core.CreatePullReqRequest) (*core.CreatePullReqResponse, error)
	GetTeamByTeamName(ctx context.Context, teamName string) (*core.GetTeamResponse, error)
	MergePullRequest(ctx context.Context, req core.MergePullReqRequest) (*core.MergePullReqResponse, error)
	ReassignPullRequest(ctx context.Context, request *core.ReassignPullReqRequest) (*core.ReassignPullReqResponse, error)
	GetReview(ctx context.Context, userID string) (*core.GetReviewResponse, error)
	UserSetIsActive(ctx context.Context, request core.SetIsActiveRequest) (*core.SetIsActiveResponse, error)
}

type Server struct {
	impl ImplInterface
}

func NewServer(impl ImplInterface) *Server {
	return &Server{
		impl: impl,
	}
}
