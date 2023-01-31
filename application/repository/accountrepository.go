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
