package core

import (
	"avito-tech/internal/app/team"
	"context"
)

// GetTeamResponse - ответ на получение команды (соответствует OpenAPI Team schema)
// OpenAPI возвращает Team schema напрямую, без обертки
type GetTeamResponse struct {
	TeamName string               `json:"team_name"`
	Members  []team.TeamMemberDTO `json:"members"`
}

func (s *Service) GetTeamByTeamName(ctx context.Context, teamName string) (*GetTeamResponse, error) {
	dto, err := s.team.GetByTeamName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	return &GetTeamResponse{
		TeamName: dto.TeamName,
		Members:  dto.Members,
	}, nil
}
