package model

type Education struct {
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	Institution string    `json:"institution"`
	Subject     string    `json:"subject"`
	Notes       []*string `json:"notes"`
}

type EducationInput struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	Institution string    `json:"institution"`
	Subject     string    `json:"subject"`
	Notes       []*string `json:"notes"`
}
