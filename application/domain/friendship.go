package domain

type Friendship struct {
	Id             string    `json:"id,omitempty"`
	Uid            string    `json:"uid,omitempty"`
	Friend         []Account `json:"friend,omitempty"`
	FriendId       string    `json:"friendId,omitempty"`
	Status         string    `json:"status,omitempty"`
	PreviousStatus string    `json:"previousStatus,omitempty"`
	ReverseStatus  string    `json:"reverseStatus,omitempty"`
	DType          []string  `json:"dgraph.type,omitempty"`
}

type Friendships struct {
	Friendships []Friendship `json:"friendships,omitempty"`
}

type CountRequest struct {
	CountRequest []Count `json:"countRequest"`
}

type IdsFriends struct {
	Ids []Id `json:"ids"`
}

type Id struct {
	Id string `json:"id"`
}

type Count struct {
	Count int `json:"count"`
}

const (
	FRIEND         = "FRIEND"
	REQUEST_TO     = "REQUEST_TO"
	REQUEST_FROM   = "REQUEST_FROM"
	BLOCKED        = "BLOCKED"
	DECLINED       = "DECLINED"
	SUBSCRIBED     = "SUBSCRIBED"
	NONE           = "NONE"
	WATCHING       = "WATCHING"
	REJECTING      = "REJECTING"
	RECOMMENDATION = "RECOMMENDATION"
)
