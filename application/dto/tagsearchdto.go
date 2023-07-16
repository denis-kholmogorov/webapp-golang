package dto

type TagSearchDto struct {
	Name string `form:"name" json:"name,omitempty"`
}
