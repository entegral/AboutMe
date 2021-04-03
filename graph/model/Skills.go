package model


type GoSkills struct {
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type GoSkillsInput struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type JSSkills struct {
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type JSSkillsInput struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type PythonSkills struct {
	Frameworks []*string `json:"frameworks"`
	Misc       []*string `json:"misc"`
}

type PythonSkillsInput struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Frameworks []*string `json:"frameworks"`
	Misc       []*string `json:"misc"`
}

type Skills struct {
	Js     *JSSkills     `json:"JS"`
	Go     *GoSkills     `json:"Go"`
	Python *PythonSkills `json:"Python"`
}
