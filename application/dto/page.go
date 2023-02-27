package dto

type PageRequest struct {
	Sort string `form:"sort" json:"sort"`
	Page int    `form:"page" json:"page"`
	Size int    `form:"size" json:"size"`
}

func (p *PageRequest) GetPage() int {
	return p.Page
}
