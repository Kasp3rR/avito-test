package core

import (
	"context"
	"time"
)

type MergePullReqRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePullReqResponse struct {
	PR MergePullReqPR `json:"pr"`
}

type MergePullReqPR struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func (s *Service) MergePullRequest(ctx context.Context, req MergePullReqRequest) (*MergePullReqResponse, error) {
	dto, err := s.pullRequest.Merge(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	response := &MergePullReqResponse{
		PR: MergePullReqPR{
			PullRequestID:     dto.PullRequestID,
			PullRequestName:   dto.PullRequestName,
			AuthorID:          dto.AuthorID,
			Status:            dto.Status,
			AssignedReviewers: dto.AssignedReviewers,
		},
	}

	if !dto.MergedAt.IsZero() {
		response.PR.MergedAt = dto.MergedAt
	}

	return response, nil
}
