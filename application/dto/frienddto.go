package dto

type FriendDto struct {
	FriendId           string `json:"friendId"`
	StatusCode         string `json:"statusCode"`
	PreviousStatusCode string `json:"previousStatusCode"`
	Rating             int64  `json:"rating"`
}
