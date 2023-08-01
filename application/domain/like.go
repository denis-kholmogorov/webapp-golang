package domain

type Like struct {
	Id           string   `json:"id,omitempty"`
	Uid          string   `json:"uid,omitempty"`
	AuthorId     string   `json:"authorId,omitempty"`
	ReactionType string   `json:"reactionType,omitempty"`
	DType        []string `json:"dgraph.type,omitempty"`
}

type LikeList struct {
	List []Like `json:"likes,omitempty"`
}

type LikeCount struct {
	ReactionType string `json:"reactionType"`
	Count        int    `json:"count"`
}

type RowMyReaction struct {
	MyReaction string `json:"myReaction"`
}
type RowReactions struct {
	Reactions []LikeCount `json:"reactions"`
}
