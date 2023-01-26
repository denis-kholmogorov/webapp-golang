package mapper

import (
	"time"
	"web/application/domain"
	"web/application/dto"
)

type Mapper interface {
	toDto(dto any) any
	toEntity(entity any) any
}

var personMapper PersonMapper
var isActivate bool

type PersonMapper struct{}

func GetPersonMapper() *PersonMapper {
	if !isActivate {
		personMapper = PersonMapper{}
	}
	return &personMapper
}

func (p *PersonMapper) toDto(person any) any {
	return person
}

func (p *PersonMapper) toEntity(person any) any {
	return person
}

func (p *PersonMapper) RegistrationToAccount(regDto dto.RegistrationDto, hashPass []byte) domain.Account {
	return domain.Account{
		FirstName: regDto.FirstName,
		LastName:  regDto.LastName,
		Email:     regDto.Email,
		Password:  string(hashPass),
		BirthDate: time.Time{},
	}

}
