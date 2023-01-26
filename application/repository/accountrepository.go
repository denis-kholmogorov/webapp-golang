package repository

import (
	"reflect"
	"web/application/domain"
)

var accountRepo *AccountRepository[*domain.Account, string]
var isInitializedAccountRepo bool
var accountType domain.Account

type AccountRepository[Entity *domain.Account, Id string] struct {
	repository *Repository
}

func GetAccountRepository() *AccountRepository[*domain.Account, string] {
	if !isInitializedAccountRepo {
		accountRepo = &AccountRepository[*domain.Account, string]{}
		accountRepo.repository = GetRepository()
		isInitializedAccountRepo = true
	}
	return accountRepo
}

func (rc *AccountRepository[Entity, Id]) FindById(id Id) (Entity, error) {
	c, err := rc.repository.FindById(accountType, string(id))
	if err != nil || c == nil {
		return nil, err
	}
	return reflect.ValueOf(c).Interface().(Entity), nil
}

func (rc *AccountRepository[Entity, Id]) FindAllBySpec(specification *Specification) ([]Entity, error) {
	rows, err := rc.repository.FindAllBySpec(accountType, specification)
	if err != nil || rows == nil {
		return nil, err
	}
	return rc.collectToListEntity(rows), nil
}

func (rc *AccountRepository[Entity, Id]) FindAll() ([]Entity, error) {
	rows, err := rc.repository.FindAll(accountType)
	if err != nil || rows == nil {
		return nil, err
	}
	return rc.collectToListEntity(rows), nil
}

func (rc *AccountRepository[Entity, Id]) collectToListEntity(collection []interface{}) []Entity {
	entities := make([]Entity, len(collection))
	for i, e := range collection {
		entities[i] = reflect.ValueOf(e).Interface().(Entity)
	}
	return entities
}
