package team

// TeamEntity - сущность для работы с БД
type TeamEntity struct {
	ID       uint64 `db:"id"`
	TeamName string `db:"team_name"`
}

// TeamMemberEntity - сущность участника команды для работы с БД
type TeamMemberEntity struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
}
