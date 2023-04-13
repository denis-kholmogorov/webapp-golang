package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"strconv"
	"web/application/domain"
	"web/application/errorhandler"
)

var settingsRepo *SettingsRepository
var isInitializedSettingsRepo bool

type SettingsRepository struct {
	conn *dgo.Dgraph
}

func GetSettingsRepository() *SettingsRepository {
	if !isInitializedSettingsRepo {
		settingsRepo = &SettingsRepository{}
		settingsRepo.conn = GetDGraphConn().connection
		isInitializedSettingsRepo = true
	}
	return settingsRepo
}

func (r SettingsRepository) GetSettings(currentUserId string) *domain.Settings {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	vars, err := txn.QueryWithVars(ctx, findSettings, variables)
	settingsList := domain.SettingsList{}
	err = json.Unmarshal(vars.Json, &settingsList)

	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("SettingsRepository:GetSettings() Error Unmarshal %s", err)})
	} else if len(settingsList.List) != 1 {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("SettingsRepository:GetSettings() Settings found more then one %s", err)})
	}
	return &settingsList.List[0]
}

func (r SettingsRepository) UpdateSettings(settings *domain.Settings, item string, value bool) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	field := map[string]string{item: strconv.FormatBool(value)}
	err := UpdateNodeFields(txn, ctx, settings.Id, field, true)

	if err != nil {
		txn.Discard(context.Background())
		panic(errorhandler.DbError{Message: fmt.Sprintf("SettingsRepository:UpdateSettings() Error AddEdges createFriendship %s", err)})
	}
}

var findSettings = `query Posts($currentUserId: string)
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
