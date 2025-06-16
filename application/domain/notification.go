package domain

import (
	"time"
)

type Notification struct {
	Id               string     `json:"id,omitempty"`
	Uid              string     `json:"uid,omitempty"`
	AuthorId         string     `json:"authorId"`
	RecipientId      string     `json:"recipientId"`
	Content          string     `json:"content"`
	NotificationType string     `json:"notificationType"`
	SentTime         *time.Time `json:"sentTime" time_format:"2006-01-02 15:04:05.99Z07:00"`
	DType            []string   `json:"dgraph.type,omitempty"`
}

func NewNotification(authorId string, recipientId string, content string, notificationType string) *Notification {
	timeNow := time.Now().UTC()
	return &Notification{
		AuthorId:         authorId,
		RecipientId:      recipientId,
		Content:          content,
		NotificationType: notificationType,
		SentTime:         &timeNow,
		DType:            []string{"Notification"},
	}
}

type EventNotification struct {
	InitiatorId      string `json:"initiatorId"`
	Content          string `json:"content"`
	NotificationType string `json:"notificationType"`
}

type NotificationList struct {
	List []Notification `json:"notifications,omitempty"`
}

func NewSettingsNotification() Settings {
	return Settings{
		Uid:                    "_:settings",
		EnablePost:             true,
		EnablePostComment:      true,
		EnableCommentComment:   true,
		EnableMessage:          true,
		EnableFriendRequest:    true,
		EnableFriendBirthday:   true,
		EnableSendEmailMessage: true,
		IsDeleted:              true,
		DType:                  []string{"Settings"},
	}
}
