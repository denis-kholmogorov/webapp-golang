package domain

type City struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	DType []string `json:"dgraph.type,omitempty"`
}
