package pullrequest

import "time"

type PullRequestEntity struct {
	PullRequestID     string     `db:"pull_request_id"`
	PullRequestName   string     `db:"pull_request_name"`
	AuthorID          string     `db:"author_id"`
	Status            string     `db:"status"`
	AssignedReviewers []string   `db:"-"`
	CreatedAt         time.Time  `db:"created_at"`
	MergedAt          *time.Time `db:"merged_at"`
}
