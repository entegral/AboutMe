package model


type Experience struct {
	StartDate        *string   `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	Title            *string   `json:"title"`
	Company          *string   `json:"company"`
	Responsibilities []*string `json:"responsibilities"`
}

type ExperienceInput struct {
	StartDate        *string   `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	Title            string    `json:"title"`
	Company          string    `json:"company"`
	Responsibilities []*string `json:"responsibilities"`
}
