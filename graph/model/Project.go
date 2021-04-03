package model


type Project struct {
	Name        string  `json:"name"`
	SourceCode  *string `json:"source_code"`
	DemoLink    *string `json:"demo_link"`
	Description *string `json:"description"`
}

type ProjectInput struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Name        string  `json:"name"`
	SourceCode  *string `json:"source_code"`
	DemoLink    *string `json:"demo_link"`
	Description *string `json:"description"`
}