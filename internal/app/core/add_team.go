package core

import (
	"avito-tech/internal/app/team"
	"context"
)

type AddTeamRequest struct {
	TeamName string               `json:"team_name"`
	Members  []team.TeamMemberDTO `json:"members"`
}

type AddTeamResponse struct {
	Team team.TeamDTO `json:"team"`
}

func (s *Service) AddTeam(ctx context.Context, req *AddTeamRequest) (*AddTeamResponse, error) {
	dto := &team.TeamDTO{
		TeamName: req.TeamName,
		Members:  req.Members,
	}
	err := s.team.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	createdTeam, err := s.team.GetByTeamName(ctx, req.TeamName)
	if err != nil {
		return nil, err
	}
	return &AddTeamResponse{Team: *createdTeam}, nil
}
