package domain

type Tag struct {
	Id    string   `json:"id,omitempty"`
	Uid   string   `json:"uid,omitempty"`
	Name  string   `json:"name"`
	DType []string `json:"dgraph.type,omitempty"`
}

type TagList struct {
	List []Tag `json:"tagList,omitempty"`
}
