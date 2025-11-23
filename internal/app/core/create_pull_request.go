package core

import (
	pullrequest "avito-tech/internal/app/pull_request"
	"context"
	"time"
)

type CreatePullReqRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePullReqResponse struct {
	PR CreatePullReqPR `json:"pr"`
}

type CreatePullReqPR struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func (s *Service) CreatePullRequestFromCreateRequest(ctx context.Context, request *CreatePullReqRequest) (*CreatePullReqResponse, error) {
	prShort := &pullrequest.PullRequestShortDTOFromHttp{
		PullRequestID:   request.PullRequestID,
		PullRequestName: request.PullRequestName,
		AuthorID:        request.AuthorID,
	}

	dto, err := s.pullRequest.Create(ctx, prShort)
	if err != nil {
		return nil, err
	}

	response := &CreatePullReqResponse{
		PR: CreatePullReqPR{
			PullRequestID:     dto.PullRequestID,
			PullRequestName:   dto.PullRequestName,
			AuthorID:          dto.AuthorID,
			Status:            dto.Status,
			AssignedReviewers: dto.AssignedReviewers,
		},
	}

	if !dto.CreatedAt.IsZero() {
		response.PR.CreatedAt = &dto.CreatedAt
	}
	if !dto.MergedAt.IsZero() {
		response.PR.MergedAt = dto.MergedAt
	}

	return response, nil
}
