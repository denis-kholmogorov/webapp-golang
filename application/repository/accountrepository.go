package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"web/application/domain"
)

var accountRepo *AccountRepository
var isInitializedAccountRepo bool

type AccountRepository struct {
	conn *dgo.Dgraph
}

func GetAccountRepository() *AccountRepository {
	if !isInitializedAccountRepo {
		accountRepo = &AccountRepository{}
		accountRepo.conn = GetDGraphConn().connection
		isInitializedAccountRepo = true
	}
	return accountRepo
}

func (r AccountRepository) Save(account *domain.Account) (*string, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	accountm, err := json.Marshal(account)
	if err != nil {
		log.Printf("AccountRepository:save() Error marhalling captcha %s", err)
		return nil, fmt.Errorf("AccountRepository:Save() Error marhalling captcha %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: accountm, CommitNow: true})
	if err != nil {
		log.Printf("AccountRepository:save() Error mutate %s", err)
		return nil, fmt.Errorf("AccountRepository:Save() Error mutate %s", err)
	}
	accountId := mutate.Uids[account.Email]
	if len(accountId) == 0 {
		log.Printf("AccountRepository:save() capthcaId not found")
		return nil, fmt.Errorf("AccountRepository:Save() capthcaId not found")
	}
	return &accountId, nil
}

func (r AccountRepository) FindByEmail(email string) ([]domain.Account, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	query := fmt.Sprintf(existEmail, email)
	vars, err := txn.Query(ctx, query)
	if err != nil {
		log.Printf("AccountRepository:existsEmail() Error query %s", err)
		return nil, fmt.Errorf("AccountRepository:existsEmail() Error query %s", err)
	}
	response := struct {
		Accounts []domain.Account `json:"accounts"`
	}{}
	err = json.Unmarshal(vars.Json, &response)

	if err != nil {
		log.Printf("AccountRepository:FindById() Error Unmarshal %s", err)
		return nil, fmt.Errorf("AccountRepository:existsEmail() Error Unmarshal %s", err)
	}
	return response.Accounts, nil

}

func (r AccountRepository) FindById(id string) (*domain.Account, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	query := fmt.Sprintf(findById, id)
	vars, err := txn.Query(ctx, query)
	if err != nil {
		log.Printf("AccountRepository:existsEmail() Error query %s", err)
		return nil, fmt.Errorf("AccountRepository:existsEmail() Error query %s", err)
	}
	response := struct {
		Accounts []domain.Account `json:"account"`
	}{}
	err = json.Unmarshal(vars.Json, &response)

	if err != nil || len(response.Accounts) != 1 {
		log.Printf("AccountRepository:FindById() Error Unmarshal %s", err)
		return nil, fmt.Errorf("AccountRepository:existsEmail() Error Unmarshal %s", err)
	}

	return &response.Accounts[0], nil

}

var existEmail = `{ accounts (func: eq(email, "%s")) {
		uid
		firstName
		lastName
		email
		password
	}
}`

var findById = `{ account (func: uid("%s")) {
		uid
		email
		firstName
		lastName
		age
		isDeleted
		isBlocked
		isOnline
		phone
		photo
		photoId
		photoName
		about
		city
		country
		birthDate
		lastOnlineTime
	}
}`
