package dto

type CaptchaDto struct {
	tableName struct{} `pg:"public.captcha"`
	Secret    string   `json:"secret"`
	Image     string   `json:"image"`
}
