package domain

type Friendship struct {
	Id             string    `json:"id,omitempty"`
	Uid            string    `json:"uid,omitempty"`
	Friend         []Account `json:"friend,omitempty"`
	Status         string    `json:"status,omitempty"`
	PreviousStatus string    `json:"previousStatus,omitempty"`
	ReverseStatus  string    `json:"reverseStatus,omitempty"`
	DType          []string  `json:"dgraph.type,omitempty"`
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
