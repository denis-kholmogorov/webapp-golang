package dto

type RegistrationDto struct {
	Email         string `json:"email"`
	Password      string `json:"password1"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	CaptchaCode   string `json:"captchaCode"`
	CaptchaSecret string `json:"captchaSecret"`
}
