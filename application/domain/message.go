package domain

import "time"

type Message struct {
	Id          string   `json:"id,omitempty"`
	Uid         string   `json:"uid,omitempty"`
	IsDeleted   bool     `json:"isDeleted"`
	IsRead      bool     `json:"isRead"`
	TimeSend    int64    `json:"timeSend"`
	Time        int64    `json:"time"` //TODO пофиксить дату с инт
	AuthorId    string   `json:"authorId"`
	RecipientId string   `json:"recipientId"`
	MessageText string   `json:"messageText"`
	DType       []string `json:"dgraph.type,omitempty"`
}

type MessageList struct {
	List []Message `json:"data,omitempty"`
}

func CreateMessage(currentUserId string, companionId string) Message {
	return Message{
		Uid:         "_:message",
		DType:       []string{"Message"},
		AuthorId:    currentUserId,
		RecipientId: companionId,
		IsRead:      true,
		TimeSend:    time.Now().UTC().Unix(),
		MessageText: "Вас добавили в чат",
	}
}
