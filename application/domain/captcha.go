package domain

import (
	"fmt"
	"time"
)

type Captcha struct {
	CaptchaId       string    `json:"uid,omitempty"`
	CaptchaCode     string    `json:"captchaCode,omitempty"`
	ExpiredTime     time.Time `json:"expiredTime" time_format:"2006-01-02 15:04:05.99Z07:00"`
	DType           []string  `json:"dgraph.type,omitempty"`
	CaptchaCodeByte []byte
}

type CaptchaList struct {
	List []Captcha `json:"captchaList,omitempty"`
}

func NewCaptchaWithCode(captchaCode string, code []byte) *Captcha {
	return &Captcha{
		CaptchaId:       fmt.Sprintf("_:%s", captchaCode),
		CaptchaCode:     captchaCode,
		ExpiredTime:     time.Now(),
		DType:           []string{"getCaptcha"},
		CaptchaCodeByte: code,
	}
}

func NewCaptcha() *Captcha {
	return &Captcha{
		DType: []string{"getCaptcha"},
	}
}
