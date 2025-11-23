package core

import (
	"avito-tech/internal/app/user"
	"context"
)

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	User user.UserDTO `json:"user"`
}

func (s *Service) UserSetIsActive(ctx context.Context, request SetIsActiveRequest) (*SetIsActiveResponse, error) {
	userDTO, err := s.user.SetIsActive(ctx, request.UserID, request.IsActive)
	if err != nil {
		return nil, err
	}
	return &SetIsActiveResponse{User: *userDTO}, nil
}
