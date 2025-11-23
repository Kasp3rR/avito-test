package team

// TeamDTO - DTO для работы с командой (используется в API)
type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (tm *TeamMemberDTO) ToEntity() TeamMemberEntity {
	return TeamMemberEntity{
		UserID:   tm.UserID,
		Username: tm.Username,
		IsActive: tm.IsActive,
	}
}

func (tm *TeamMemberDTO) FromEntity(entity *TeamMemberEntity) {
	tm.UserID = entity.UserID
	tm.Username = entity.Username
	tm.IsActive = entity.IsActive
}

func ToTeamMeberEntities(tms []TeamMemberDTO) []TeamMemberEntity {
	entities := make([]TeamMemberEntity, len(tms))
	for i, v := range tms {
		entities[i] = v.ToEntity()
	}
	return entities
}

func FromTeamMeberEntities(tme []TeamMemberEntity) []TeamMemberDTO {
	dto := make([]TeamMemberDTO, len(tme))
	for i, v := range tme {
		var m TeamMemberDTO
		m.FromEntity(&v)
		dto[i] = m
	}
	return dto
}
