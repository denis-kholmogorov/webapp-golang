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

var captchaRepo *CaptchaRepository
var isInitializedCaptchaRepo bool

type CaptchaRepository struct {
	conn *dgo.Dgraph
}

func GetCaptchaRepository() *CaptchaRepository {
	if !isInitializedCaptchaRepo {
		captchaRepo = &CaptchaRepository{}
		captchaRepo.conn = GetDGraphConn().connection
		isInitializedCaptchaRepo = true
	}
	return captchaRepo
}

func (r CaptchaRepository) FindById(captchaId string) (*domain.Captcha, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	query := fmt.Sprintf(getById, captchaId)
	vars, err := txn.Query(ctx, query)
	if err != nil {
		log.Printf("CaptchaRepository:FindById() Error query %s", err)
		return nil, fmt.Errorf("CaptchaRepository:FindById() Error query %s", err)
	}

	captchaList := domain.CaptchaList{}

	err = json.Unmarshal(vars.Json, &captchaList)

	if err != nil {
		log.Printf("CaptchaRepository:FindById() Error Unmarshal %s", err)
		return nil, fmt.Errorf("CaptchaRepository:FindById() Error Unmarshal %s", err)
	}

	if len(captchaList.List) != 1 {
		log.Printf("CaptchaRepository:FindById() Captha found more then one %s", err)
		return nil, fmt.Errorf("CaptchaRepository:FindById() Captha found more then one %s", err)
	}
	return &captchaList.List[0], nil
}

func (r CaptchaRepository) Save(captcha *domain.Captcha) (*string, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	captcham, err := json.Marshal(captcha)
	if err != nil {
		log.Printf("CaptchaRepository:save() Error marhalling captcha %s", err)
		return nil, fmt.Errorf("CaptchaRepository:Create() Error marhalling captcha %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: captcham, CommitNow: true})
	if err != nil {
		log.Printf("CaptchaRepository:save() Error mutate %s", err)
		return nil, fmt.Errorf("CaptchaRepository:Create() Error mutate %s", err)
	}
	captchaId := mutate.Uids[captcha.CaptchaCode]
	if len(captchaId) == 0 {
		log.Printf("CaptchaRepository:save() capthcaId not found")
		return nil, fmt.Errorf("CaptchaRepository:Create() capthcaId not found")
	}
	return &captchaId, nil
}

var deleteById = `{
	delete {
		<%uid * * .>
	}	
}`

var getById = `{ captchaList (func: uid(%s)) {
uid
captchaCode 
expiredTime
}
}`
