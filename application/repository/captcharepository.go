package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"web/application/domain"
	"web/application/errorhandler"
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

func (r CaptchaRepository) FindById(captchaId string) *domain.Captcha {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	query := fmt.Sprintf(getById, captchaId)
	vars, err := txn.Query(ctx, query)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("CaptchaRepository:FindById() Error query %s", err)})
	}

	captchaList := domain.CaptchaList{}
	err = json.Unmarshal(vars.Json, &captchaList)

	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("CaptchaRepository:FindById() Error Unmarshal %s", err)})
	} else if len(captchaList.List) != 1 {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("CaptchaRepository:FindById() Captha found more then one %s", err)})
	}

	return &captchaList.List[0]
}

func (r CaptchaRepository) Save(captcha *domain.Captcha) string {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	captcham, err := json.Marshal(captcha)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("CaptchaRepository:Save() Error Unmarshal %s", err)})
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: captcham, CommitNow: true})
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("CaptchaRepository:save() Error mutate %s", err)})
	}
	captchaId := mutate.Uids[captcha.CaptchaCode]
	if len(captchaId) == 0 {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("CaptchaRepository:save() capthcaId not saved %s", err)})
	}
	return captchaId
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
