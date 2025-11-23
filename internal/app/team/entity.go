package team

type TeamEntity struct {
	ID       uint64 `db:"id"`
	TeamName string `db:"team_name"`
}

type TeamMemberEntity struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
}
