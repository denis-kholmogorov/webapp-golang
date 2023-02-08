package domain

type Country struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Cities []City   `json:"cities"`
	DType  []string `json:"dgraph.type,omitempty"`
}

type CountriesList struct {
	List []Country `json:"countriesList,omitempty"`
}
