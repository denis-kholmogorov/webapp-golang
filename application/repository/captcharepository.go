package repository

import (
	"reflect"
	"web/application/domain"
)

var captchaRepo *CaptchaRepository[*domain.Captcha, string]
var isInitializedCaptchaRepo bool

type CaptchaRepository[Entity *domain.Captcha, Id string] struct {
	repository *Repository
}

func GetCaptchaRepository() *CaptchaRepository[*domain.Captcha, string] {
	if !isInitializedCaptchaRepo {
		captchaRepo = &CaptchaRepository[*domain.Captcha, string]{}
		captchaRepo.repository = GetRepository()
		isInitializedCaptchaRepo = true
	}
	return captchaRepo
}

func (rc *CaptchaRepository[Entity, Id]) FindById(domainType Entity, id Id) (Entity, error) {
	c, err := rc.repository.FindById(*domainType, string(id))
	if err != nil || c == nil {
		return nil, err
	}
	return reflect.ValueOf(c).Interface().(Entity), nil
}
