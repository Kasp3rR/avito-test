package user

type UserEntity struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	TeamID   uint64 `db:"team_id"`
	TeamName string `db:"-"`
	IsActive bool   `db:"is_active"`
}
