package domain

type Like struct {
	Id       string   `json:"id,omitempty"`
	Uid      string   `json:"uid,omitempty"`
	AuthorId string   `json:"authorId,omitempty"`
	DType    []string `json:"dgraph.type,omitempty"`
}

type LikeList struct {
	List []Like `json:"likeList,omitempty"`
}

type LikeExist struct {
	Exists []LikeCount `json:"exists,omitempty"`
}

type LikeCount struct {
	Count int `json:"count,omitempty"`
}
