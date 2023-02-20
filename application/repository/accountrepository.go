package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"math"
	"strconv"
	"time"
	"web/application/domain"
	"web/application/dto"
	"web/application/utils"
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
		log.Printf("AccountRepository:save() Error marhalling account %s", err)
		return nil, fmt.Errorf("AccountRepository:Create() Error marhalling account %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: accountm, CommitNow: true})
	if err != nil {
		log.Printf("AccountRepository:save() Error mutate %s", err)
		return nil, fmt.Errorf("AccountRepository:Create() Error mutate %s", err)
	}
	accountId := mutate.Uids[account.Email]
	if len(accountId) == 0 {
		log.Printf("AccountRepository:save() capthcaId not found")
		return nil, fmt.Errorf("AccountRepository:Create() capthcaId not found")
	}
	return &accountId, nil
}

func (r AccountRepository) FindByEmail(email string) ([]domain.Account, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$email"] = email
	vars, err := txn.QueryWithVars(ctx, existEmail, variables)
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
	variables := make(map[string]string)
	variables["$id"] = id
	vars, err := txn.QueryWithVars(ctx, findById, variables)
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

func (r AccountRepository) Update(account *domain.Account) (*string, error) {
	timeNow := time.Now().UTC()
	account.DType = []string{"Account"}
	account.UpdatedOn = &timeNow
	ctx := context.Background()
	txn := r.conn.NewTxn()
	accountm, err := json.Marshal(account)
	if err != nil {
		log.Printf("AccountRepository:Update() Error marhalling account %s", err)
		return nil, fmt.Errorf("AccountRepository:Update() Error marhalling account %s", err)
	}
	_, err = txn.Mutate(ctx, &api.Mutation{SetJson: accountm, CommitNow: true})
	if err != nil {
		log.Printf("AccountRepository:Update() Error mutate %s", err)
		return nil, fmt.Errorf("AccountRepository:Update() Error mutate %s", err)
	}
	return &account.Uid, nil
}

func (r AccountRepository) FindAll(dto dto.AccountSearchDto) (*domain.PageResponse, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$search"] = fmt.Sprintf("/.*%s.*/", dto.Author)
	variables["$first"] = strconv.Itoa(dto.Size)
	variables["$offset"] = strconv.Itoa(dto.Size * utils.GetPageNumber(&dto))

	var vars *api.Response
	var err error

	if len(dto.Author) > 2 {
		vars, err = txn.QueryWithVars(ctx, findByAuthor, variables)
	} else if len(dto.Author) > 0 && len(dto.Author) <= 2 {
		log.Printf("GeoRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("query should be more to %s words", "2")
	}

	response := domain.PageResponse{}
	err = json.Unmarshal(vars.Json, &response)
	if len(response.Content) > 0 {
		response.Size = dto.Size
		response.TotalElement = response.Count[0].TotalElement
		response.TotalPages = int(math.Ceil(float64(response.TotalElement) / float64(dto.Size)))
		response.Number = dto.Page + 1
	}
	return &response, nil

}

var existEmail = `query AccountByEmail($email: string)
{ accounts (func: eq(email, $email)) {
		id:uid
		firstName
		lastName
		email
		password
	}
}`

var findById = `query AccountById($id: string)
{ account (func: uid($id)) {
		id:uid
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

var findByAuthor = `query searchAuthor($search: string, $first: int, $offset: int){
var(func: regexp(firstName, $search), first: 10, offset: 0)  @filter(eq(dgraph.type, Account) and eq(isDeleted, false)){
    A as uid 
  }  
    var(func: regexp(lastName, $search), first: 10, offset: 0)  @filter(eq(dgraph.type, Account) and eq(isDeleted, false)){
    B as uid 
  }
	content(func: uid(A,B), orderdesc: firstName, first: 10, offset: 0)  {
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
	}
	count(func: type(Account)) @filter(eq(isDeleted,false)){
		totalElement: count(firstName)
	}
}`

//
//var findByAuthor = `query searchAuthor($search: string, $first: int, $offset: int)
//{ content(func: regexp(firstName, $search),first: $first, offset: $offset) @filter(eq(dgraph.type, Account) and eq(isDeleted, false))  {
//		firstName
//		lastName
//		age
//		isDeleted
//		isBlocked
//		isOnline
//		phone
//		photo
//		photoId
//		photoName
//}
//totalElement(func: type(Account)) @filter(eq(isDeleted,false)){
//				count: count(firstName)
//      }
//}`
