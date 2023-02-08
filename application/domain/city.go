package domain

type City struct {
	Title string   `json:"title"`
	DType []string `json:"dgraph.type,omitempty"`
}
