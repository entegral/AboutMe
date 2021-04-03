package model


type ContactInfo struct {
	Email    string  `json:"email"`
	LinkedIn *string `json:"linkedIn"`
	Github   *string `json:"github"`
}

