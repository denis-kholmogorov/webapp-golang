package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"strconv"
	"strings"
	"web/application/domain"
	"web/application/errorhandler"
)

var settingsRepo *NotificationRepository
var isInitializedSettingsRepo bool

type NotificationRepository struct {
	conn *dgo.Dgraph
}

func GetSettingsRepository() *NotificationRepository {
	if !isInitializedSettingsRepo {
		settingsRepo = &NotificationRepository{}
		settingsRepo.conn = GetDGraphConn().connection
		isInitializedSettingsRepo = true
	}
	return settingsRepo
}

func (r NotificationRepository) Create(event *domain.EventNotification) error {
	//ctx := context.Background()
	//txn := r.conn.NewTxn()
	////field := map[string]string{item: strconv.FormatBool(value)}
	//err := UpdateNodeFields(txn, ctx, settings.Id, field, true)
	//
	//if err != nil {
	//	txn.Discard(context.Background())
	//	panic(errorhandler.DbError{Message: fmt.Sprintf("NotificationRepository:UpdateSettings() Error AddEdges createFriendship %s", err)})
	//}
	return nil
}

func (r NotificationRepository) GetAll(currentUserId string) []domain.Notification {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	vars, err := txn.QueryWithVars(ctx, getAll, variables)
	notifyList := domain.NotificationList{}
	err = json.Unmarshal(vars.Json, &notifyList)

	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("NotificationRepository:GetSettings() Settings found more then one %s", err)})
	}
	return notifyList.List
}

func (r NotificationRepository) GetSettings(currentUserId string) *domain.Settings {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	vars, err := txn.QueryWithVars(ctx, getSettings, variables)
	settingsList := domain.SettingsList{}
	err = json.Unmarshal(vars.Json, &settingsList)

	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("NotificationRepository:GetSettings() Error Unmarshal %s", err)})
	} else if len(settingsList.List) != 1 {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("NotificationRepository:GetSettings() Settings found more then one %s", err)})
	}
	return &settingsList.List[0]
}

func (r NotificationRepository) UpdateSettings(settings *domain.Settings, item string, value bool) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	field := map[string]string{item: strconv.FormatBool(value)}
	err := UpdateNodeFields(txn, ctx, settings.Id, field, true)

	if err != nil {
		txn.Discard(context.Background())
		panic(errorhandler.DbError{Message: fmt.Sprintf("NotificationRepository:UpdateSettings() Error AddEdges createFriendship %s", err)})
	}
}

func (r NotificationRepository) CreateSettings(settings *domain.Settings, item string, value bool) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	field := map[string]string{item: strconv.FormatBool(value)}
	err := UpdateNodeFields(txn, ctx, settings.Id, field, true)

	if err != nil {
		txn.Discard(context.Background())
		panic(errorhandler.DbError{Message: fmt.Sprintf("NotificationRepository:UpdateSettings() Error AddEdges createFriendship %s", err)})
	}
}

func (r NotificationRepository) SaveAll(notifications []*domain.Notification) (*domain.NotificationList, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	notificationsm, err := json.Marshal(notifications)
	if err != nil {
		log.Printf("NotificationRepository:Save() Error marhalling account %s", err)
		return nil, err
	}
	mu, err := txn.Mutate(ctx, &api.Mutation{SetJson: notificationsm, CommitNow: true})
	if err != nil {
		log.Printf("NotificationRepository:save() Error mutate %s", err)
		return nil, err
	}

	ids := make([]string, 0)
	for _, id := range mu.GetUids() {
		ids = append(ids, id)
	}
	txn = r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$ids"] = strings.Join(ids, ",")
	vars, err := txn.QueryWithVars(ctx, findByIds, variables)
	if err != nil {
		log.Printf("NotificationRepository:save() Error QueryWithVars %s", err)
		return nil, err
	}
	response := domain.NotificationList{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("NotificationRepository:FindAll() Error Unmarshal %s", err)
		return nil, fmt.Errorf("NotificationRepository:FindAll() Error Unmarshal %s", err)
	}
	return &response, nil
}

var getAll = `query getNotification($currentUserId: string)
{
  notifications(func: type(Notification)) @filter(eq(recipientId,$currentUserId)) {
		id: uid
        uid
  		authorId
	    recipientId
	    content
	    notificationType
	    sentTime
  }
}
`

var getSettings = `query getNotificationSettings($currentUserId: string)
{
  settingsList(func: uid($currentUserId)) @normalize {
    settings{
		id: uid
        uid: uid
  		enableCommentComment: enableCommentComment
    	enablePost: enablePost
		enablePostComment: enablePostComment
		enableMessage: enableMessage
		enableFriendRequest: enableFriendRequest
		enableFriendBirthday: enableFriendBirthday
		enableSendEmailMessage: enableSendEmailMessage
    }
  }
}
`

var findByIds = `query getNotificationByIds($ids: string)
{
  notifications(func: uid($ids)) @normalize {
		id: uid
        uid: uid
  		authorId: authorId
	    recipientId: recipientId
	    content: content
	    notificationType: notificationType
	    sentTime: sentTime
  }
}
`
