package domain

type City struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	CountryId string   `json:"countryId"`
	DType     []string `json:"dgraph.type,omitempty"`
}
