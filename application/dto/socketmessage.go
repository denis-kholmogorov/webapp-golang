package dto

import "web/application/domain"

type SocketDto[T any] struct {
	Type      string `json:"type"`
	AccountId string `json:"accountId"`
	Data      *T     `json:"data"`
}

func NewMessageSocketDto(data *domain.Message) *SocketDto[domain.Message] {
	s := SocketDto[domain.Message]{}
	s.Data = data
	s.AccountId = data.AuthorId
	s.Type = "MESSAGE"
	return &s
}
