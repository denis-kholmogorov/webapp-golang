package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"web/application/domain"
	"web/application/dto"
)

var dialogRepo *DialogRepository
var isInitializedDialogRepo bool

type DialogRepository struct {
	conn *dgo.Dgraph
}

func GetDialogRepository() *DialogRepository {
	if !isInitializedDialogRepo {
		dialogRepo = &DialogRepository{}
		dialogRepo.conn = GetDGraphConn().connection
		isInitializedDialogRepo = true
	}
	return dialogRepo
}

func (r DialogRepository) GetDialogs(currentUserId string, page dto.PageRequestOld) (*dto.PageResponse[domain.Dialog], error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$first"] = "100"
	variables["$offset"] = "0"
	vars, err = txn.QueryWithVars(ctx, getDialogs, variables)

	if err != nil {
		log.Printf("DialogRepository:GetDialogs() Error query %s", err)
		return nil, fmt.Errorf("DialogRepository:GetDialogs() Error query %s", err)
	}

	response := dto.PageResponse[domain.Dialog]{}
	err = json.Unmarshal(vars.Json, &response)

	if err != nil {
		log.Printf("DialogRepository:GetDialogs() Error Unmarshal %s", err)
		return nil, fmt.Errorf("DialogRepository:GetDialogs() Error Unmarshal %s", err)
	}
	response.SetPage(100, 0)
	return &response, nil
}

func (r DialogRepository) GetMessages(currentUserId string, page dto.PageRequestOld) (interface{}, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$companionId"] = page.CompanionId
	vars, err = txn.QueryWithVars(ctx, existDialog, variables)

	if err != nil {
		txn.Discard(ctx)
		log.Printf("DialogRepository:GetMessages() Error query %s", err)
		return nil, fmt.Errorf("DialogRepository:GetMessages() Error query %s", err)
	}
	exist := domain.DialogList{}
	err = json.Unmarshal(vars.Json, &exist)
	if err != nil {
		txn.Discard(ctx)
		log.Printf("DialogRepository:GetDialogs() Error Unmarshal %s", err)
		return nil, fmt.Errorf("DialogRepository:GetMessages() Error Unmarshal %s", err)
	}

	if len(exist.List) == 0 {
		dialog, err := createDialog(ctx, txn, currentUserId, page.CompanionId)
		return dialog, err
	}
	variablesDialog := make(map[string]string)
	variablesDialog["$dialogId"] = exist.List[0].Uid
	vars, err = txn.QueryWithVars(ctx, getMessageByDialog, variablesDialog)

	if err != nil {
		txn.Discard(ctx)
		log.Printf("DialogRepository:GetMessges() Error query %s", err)
		return nil, fmt.Errorf("DialogRepository:GetMessges() Error query %s", err)
	}
	messages := domain.MessageList{}
	err = json.Unmarshal(vars.Json, &messages)
	if err != nil {
		txn.Discard(ctx)
		log.Printf("DialogRepository:GetDialogs() Error Unmarshal %s", err)
		return nil, fmt.Errorf("DialogRepository:GetDialogs() Error Unmarshal %s", err)
	}
	txn.Commit(ctx)
	return &messages, nil
}

func createDialog(ctx context.Context, txn *dgo.Txn, currentUserId string, companionId string) (*domain.MessageList, error) {
	dialog := domain.Dialog{}
	dialog.Uid = "_:dialog"
	dialog.DType = []string{"Dialog"}
	dialog.Messages = []domain.Message{domain.CreateMessage(currentUserId, companionId)}
	marshal, err := json.Marshal(dialog)

	if err != nil {
		txn.Discard(ctx)
		log.Printf("DialogRepository:createDialog() Error marhalling post %s", err)
		return nil, fmt.Errorf("DialogRepository:createDialog() Error marhalling post %s", err)
	}

	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: marshal})

	if err != nil {
		txn.Discard(ctx)
		log.Printf("DialogRepository:createDialog() Error mutate %s", err)
		return nil, fmt.Errorf("DialogRepository:createDialog() Error mutate %s", err)
	}

	edges := []dto.Edge{
		{mutate.GetUids()["dialog"], "conversationPartner1", currentUserId}, // текущий добавляет дружбу TO // другу добавляет дружбу FROM
		{mutate.GetUids()["dialog"], "conversationPartner2", companionId},   // дружба друга добавляет текущего
	}
	err = AddEdges(txn, ctx, edges, true)
	if err != nil {
		txn.Discard(ctx)
		log.Printf("DialogRepository:createDialog() Error mutate %s", err)
		return nil, fmt.Errorf("DialogRepository:createDialog() Error mutate %s", err)
	}
	dialog.Messages[0].Id = mutate.Uids["message"]
	dialog.Messages[0].Time = dialog.Messages[0].TimeSend

	return &domain.MessageList{List: dialog.Messages}, nil
}

var getDialogs = `query Posts($currentUserId: string, $first: int, $offset: int)
{
	var(func: type(Dialog)) @filter(uid_in(participantOne,$currentUserId) or uid_in(participantTwo,$currentUserId)) {
		A as uid
	}
	{
		content(func: uid(A), first:100, offset:0) {
    		id:uid
    		unreadCount
    		conversationPartner1: participantOne @filter(not uid($currentUserId)){
      			id:uid
				firstName
      			lastName
      			photo
    		}
    		conversationPartner2: participantTwo @filter(not uid($currentUserId)){
      			id:uid
				firstName
				lastName
      			photo
    		}
  		}
		count(func: uid(A)){
			totalElement:count(uid)
		}
	}
}`

var existDialog = `query HasDialog($currentUserId: string, $companionId: string)
{
	dialogList(func: type(Dialog)) @filter((uid_in(conversationPartner2,$companionId) and uid_in(conversationPartner1,$currentUserId)) or
		(uid_in(conversationPartner2,$currentUserId) and uid_in(conversationPartner1,$companionId))) {
   		uid
	} 
}`

var getMessageByDialog = `query Messages($dialogId: string)
{
	data(func: uid($dialogId)) @normalize  {
    	messages {
   			id: uid
			time: timeSend
			messageText: messageText
        	authorId: authorId
			recipientId: recipientId
			isRead: isRead
		} 
    }
}`
