package team

import "context"

type Repo interface {
	create(ctx context.Context, teamName string, members []*TeamMemberEntity) error
	getByName(ctx context.Context, teamName string) (*TeamEntity, []TeamMemberEntity, error)
}

type Team struct {
	repo Repo
}

func NewTeam(repo Repo) *Team {
	return &Team{
		repo: repo,
	}
}

func (t *Team) GetByTeamName(ctx context.Context, teamName string) (*TeamDTO, error) {
	entity, entities, err := t.repo.getByName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	members := FromTeamMeberEntities(entities)
	var dto TeamDTO
	dto.TeamName = entity.TeamName
	dto.Members = members
	return &dto, nil
}

func (t *Team) Create(ctx context.Context, dto *TeamDTO) error {
	entities := ToTeamMeberEntities(dto.Members)
	members := make([]*TeamMemberEntity, len(entities))
	for i := range entities {
		members[i] = &entities[i]
	}
	err := t.repo.create(ctx, dto.TeamName, members)
	return err
}
