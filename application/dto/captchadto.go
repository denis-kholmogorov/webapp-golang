package dto

type CaptchaDto struct {
	Secret string `json:"secret"`
	Image  string `json:"image"`
}
