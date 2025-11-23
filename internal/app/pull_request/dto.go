package pullrequest

import "time"

type PullRequestShortDTO struct {
	PullRequestID   string `db:"pull_request_id"`
	PullRequestName string `db:"pull_request_name"`
	AuthorID        string `db:"author_id"`
	Status          string `db:"status"`
}

type PullRequestShortDTOFromHttp struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func (pr *PullRequestShortDTOFromHttp) MapToModel(prs *PullRequestShortDTO) {
	pr.PullRequestID = prs.PullRequestID
	pr.PullRequestName = prs.PullRequestName
	pr.AuthorID = prs.AuthorID
	pr.Status = prs.Status
}

func (pr *PullRequestShortDTOFromHttp) MapToPREntity() *PullRequestEntity {
	return &PullRequestEntity{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
	}
}

func MapToModelsShort(entities []PullRequestShortDTO) []*PullRequestShortDTOFromHttp {
	answer := make([]*PullRequestShortDTOFromHttp, len(entities))
	for i, v := range entities {
		var pr PullRequestShortDTOFromHttp
		pr.MapToModel(&v)
		answer[i] = &pr
	}
	return answer
}

type PullRequestDTOFromHttp struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"created_at"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func (pr *PullRequestDTOFromHttp) MapToModel() *PullRequestEntity {
	return &PullRequestEntity{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
	}
}

func (pr *PullRequestDTOFromHttp) MapFromModel(entity *PullRequestEntity) {
	pr.PullRequestID = entity.PullRequestID
	pr.PullRequestName = entity.PullRequestName
	pr.AuthorID = entity.AuthorID
	pr.Status = entity.Status
	pr.AssignedReviewers = entity.AssignedReviewers
	pr.CreatedAt = entity.CreatedAt
	pr.MergedAt = entity.MergedAt
}
