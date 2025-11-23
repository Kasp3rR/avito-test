package core

import (
	"avito-tech/internal/app/user"
	"context"
)

// SetIsActiveRequest - запрос на изменение активности пользователя (соответствует OpenAPI)
type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// SetIsActiveResponse - ответ на изменение активности (обернут в user)
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
