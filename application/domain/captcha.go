package domain

import "time"

type Captcha struct {
	tableName   struct{}  `pg:"public.captcha"`
	CaptchaId   int       `pg:"id"`
	CaptchaCode string    `pg:"captcha_code"`
	ExpiredTime time.Time `pg:"expired_time" time_format:"2006-01-02 15:04:05.99Z07:00"`
}
