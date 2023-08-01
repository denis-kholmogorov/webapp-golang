package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"strconv"
	"time"
	"web/application/domain"
	"web/application/dto"
	"web/application/errorhandler"
	"web/application/utils"
)

const regexAuthor = "/(?i).*%s.*/"

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

func (r AccountRepository) Save(account *domain.Account) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	accountm, err := json.Marshal(account)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("AccountRepository:Save() Error marhalling account %s", err)})
	}
	mu, err := txn.Mutate(ctx, &api.Mutation{SetJson: accountm, CommitNow: true})
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("AccountRepository:save() Error mutate %s", err)})
	}
	mu.GetUids()
}

func (r AccountRepository) FindByEmail(email string) []domain.Account {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$email"] = email
	vars, err := txn.QueryWithVars(ctx, existEmail, variables)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("AccountRepository:FindByEmail() Error query %s", err)})
	}
	response := struct {
		Accounts []domain.Account `json:"accounts"`
	}{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("AccountRepository:FindByEmail() Error Unmarshal %s", err)})
	}
	return response.Accounts

}

func (r AccountRepository) FindById(id string) *domain.Account {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$id"] = id
	vars, err := txn.QueryWithVars(ctx, findById, variables)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("AccountRepository:existsEmail() Error query %s", err)})
	}
	response := struct {
		Accounts []domain.Account `json:"account"`
	}{}
	err = json.Unmarshal(vars.Json, &response)

	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("AccountRepository:FindById() Error Unmarshal %s", err)})
	}
	if len(response.Accounts) != 1 {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("AccountRepository:FindById() Error not found by id %s", err)})
	}

	return &response.Accounts[0]

}

func (r AccountRepository) Update(account *domain.Account) string {
	timeNow := time.Now().UTC()
	account.DType = []string{"Account"}
	account.UpdatedOn = &timeNow
	ctx := context.Background()
	txn := r.conn.NewTxn()
	accountm, err := json.Marshal(account)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("AccountRepository:UpdateSettings() Error marhalling account %s", err)})
	}
	_, err = txn.Mutate(ctx, &api.Mutation{SetJson: accountm, CommitNow: true})
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("AccountRepository:UpdateSettings() Error mutate %s", err)})
	}
	return account.Id
}

func (r AccountRepository) FindAll(searchDto dto.AccountSearchDto, currentId string) *dto.PageResponse[domain.Account] {
	variables := make(map[string]string)
	variables["$first"] = strconv.Itoa(searchDto.Size)
	variables["$offset"] = strconv.Itoa(searchDto.Size * utils.GetPageNumber(&searchDto))
	variables["$currentId"] = currentId

	query := createSearchQuery()

	var vars *api.Response
	var err error

	txn := r.conn.NewReadOnlyTxn()

	if len(searchDto.Author) > 0 {
		variables["$search"] = fmt.Sprintf(regexAuthor, searchDto.Author)
		vars, err = txn.QueryWithVars(context.Background(), findByAuthor, variables)
	} else {
		addFilters(searchDto, variables, query)
		vars, err = txn.QueryWithVars(context.Background(), fmt.Sprintf(findByParams, query[0], query[1]), variables)
	}

	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("AccountRepository:FindAll() Error query %s", err)})
	}

	response := dto.PageResponse[domain.Account]{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("AccountRepository:FindAll() Error Unmarshal %s", err)})
	}
	response.SetPage(searchDto.Size, searchDto.Page)
	return &response
}

func createSearchQuery() map[int]string {
	m := make(map[int]string)
	m[0] = ""
	m[1] = ""
	return m
}

func addFilters(searchDto dto.AccountSearchDto, variables map[string]string, query map[int]string) {
	variables["$firstName"] = fmt.Sprintf(regexAuthor, searchDto.FirstName)
	variables["$lastName"] = fmt.Sprintf(regexAuthor, searchDto.LastName)

	if searchDto.Country != "" {
		variables["$country"] = searchDto.Country
		query[0] = query[0] + country
		query[1] = query[1] + andCountry
	}
	if searchDto.City != "" {
		variables["$city"] = searchDto.City
		query[0] = query[0] + city
		query[1] = query[1] + andCity
	}
	if searchDto.AgeFrom != 0 {
		now := time.Now()
		year := now.Year() - searchDto.AgeFrom
		variables["$ageFrom"] = time.Date(year, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02T03:04:05Z")
		query[0] = query[0] + ageFrom
		query[1] = query[1] + andAgeFrom
	}
	if searchDto.AgeTo != 0 {
		now := time.Now()
		year := now.Year() - searchDto.AgeTo
		variables["$ageTo"] = time.Date(year, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02T03:04:05Z")
		query[0] = query[0] + ageTo
		query[1] = query[1] + andAgeTo
	}
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
        emojiStatus
        profileCover
	}
}`

var findByParams = `query searchAuthor($firstName: string, $lastName: string, $currentId: string%s, $first: int, $offset: int)
{
	var(func: type(Account)) @filter(not uid($currentId) and regexp(firstName, $firstName) and regexp(lastName, $lastName)%s ){
	A as uid
}
	content(func: uid(A), orderdesc: firstName, first: $first, offset: $offset)  {
		id: uid
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
	count(func: uid(A)){
		totalElement: count(uid)
	}
}`

var findByAuthor = `query searchAuthor($search: string, $currentId: string, $first: int, $offset: int)
{
	var(func: type(Account)) @filter(not uid($currentId) and (regexp(firstName, $search) or regexp(lastName, $search))){
	A as uid
}
	content(func: uid(A), orderdesc: firstName, first: $first, offset: $offset)  {
		id: uid
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
	count(func: uid(A)){
		totalElement: count(uid)
	}
}`

var country = ", $country: string"
var city = ", $city: string"
var ageFrom = ", $ageFrom: string"
var ageTo = ", $ageTo: string"
var andCountry = "and eq(country, $country)"
var andCity = "and eq(city, $city)"
var andAgeFrom = "and lt(birthDate, $ageFrom)"
var andAgeTo = "and gt(birthDate, $ageTo)"
