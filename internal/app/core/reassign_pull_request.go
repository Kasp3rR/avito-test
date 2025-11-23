package core

import (
	"context"
	"time"
)

// ReassignPullReqRequest - запрос на переназначение ревьювера (соответствует OpenAPI)
type ReassignPullReqRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

// ReassignPullReqResponse - ответ на переназначение (обернут в pr, с replaced_by и camelCase для дат)
type ReassignPullReqResponse struct {
	PR         ReassignPullReqPR `json:"pr"`
	ReplacedBy string            `json:"replaced_by"`
}

type ReassignPullReqPR struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func (s *Service) ReassignPullRequest(ctx context.Context, request *ReassignPullReqRequest) (*ReassignPullReqResponse, error) {
	dto, newID, err := s.pullRequest.Reassign(ctx, request.PullRequestID, request.OldUserID)
	if err != nil {
		return nil, err
	}

	response := &ReassignPullReqResponse{
		ReplacedBy: newID,
		PR: ReassignPullReqPR{
			PullRequestID:     dto.PullRequestID,
			PullRequestName:   dto.PullRequestName,
			AuthorID:          dto.AuthorID,
			Status:            dto.Status,
			AssignedReviewers: dto.AssignedReviewers,
		},
	}

	// Преобразуем даты из snake_case в camelCase
	if !dto.CreatedAt.IsZero() {
		response.PR.CreatedAt = &dto.CreatedAt
	}

	return response, nil
}
