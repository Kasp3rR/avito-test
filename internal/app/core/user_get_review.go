package core

import (
	pullrequest "avito-tech/internal/app/pull_request"
	"context"
)

// GetReviewResponse - ответ на получение PR пользователя (соответствует OpenAPI)
type GetReviewResponse struct {
	UserID       string                                     `json:"user_id"`
	PullRequests []*pullrequest.PullRequestShortDTOFromHttp `json:"pull_requests"`
}

func (s *Service) GetReview(ctx context.Context, userID string) (*GetReviewResponse, error) {
	dto, err := s.user.GetReview(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetReviewResponse{
		UserID:       userID,
		PullRequests: dto,
	}, nil
}
