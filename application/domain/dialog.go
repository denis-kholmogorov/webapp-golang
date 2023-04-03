package domain

type Dialog struct {
	Id                   string    `json:"id,omitempty"`
	Uid                  string    `json:"uid,omitempty"`
	IsDeleted            bool      `json:"isDeleted"`
	UnreadCount          int       `json:"unreadCount"`
	ConversationPartner1 Account   `json:"conversationPartner1,omitempty"`
	ConversationPartner2 Account   `json:"conversationPartner2,omitempty"`
	Messages             []Message `json:"messages"`
	LastMessage          Message   `json:"lastMessage"`
	DType                []string  `json:"dgraph.type,omitempty"`
}

type DialogList struct {
	List []Dialog `json:"dialogList,omitempty"`
}

type DialogsCount struct {
	List []Count `json:"dialogsCount,omitempty"`
}
