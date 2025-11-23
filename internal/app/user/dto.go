package user

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func (u *UserDTO) MapToModel() *UserEntity {
	return &UserEntity{
		UserID:   u.UserID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func (u *UserDTO) MapFromModel(entity *UserEntity) {
	u.UserID = entity.UserID
	u.Username = entity.Username
	u.TeamName = entity.TeamName
	u.IsActive = entity.IsActive
}

func MapFromModels(entities []*UserEntity) []*UserDTO {
	dto := make([]*UserDTO, len(entities))
	for i, v := range entities {
		var a UserDTO
		a.MapFromModel(v)
		dto[i] = &a
	}
	return dto
}
