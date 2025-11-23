package core

import (
	pullrequest "avito-tech/internal/app/pull_request"
	"context"
	"time"
)

// CreatePullReqRequest - запрос на создание PR (соответствует OpenAPI)
type CreatePullReqRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

// CreatePullReqResponse - ответ на создание PR (обернут в pr, с camelCase для дат)
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
	// Преобразуем Request в DTO для сервисного слоя
	prShort := &pullrequest.PullRequestShortDTOFromHttp{
		PullRequestID:   request.PullRequestID,
		PullRequestName: request.PullRequestName,
		AuthorID:        request.AuthorID,
	}

	dto, err := s.pullRequest.Create(ctx, prShort)
	if err != nil {
		return nil, err
	}

	// Преобразуем DTO в Response с правильным форматом
	response := &CreatePullReqResponse{
		PR: CreatePullReqPR{
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
	if !dto.MergedAt.IsZero() {
		response.PR.MergedAt = dto.MergedAt
	}

	return response, nil
}
