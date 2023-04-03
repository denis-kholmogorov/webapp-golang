package dto

type PageRequest struct {
	Sort string `form:"sort" json:"sort"`
	Page int    `form:"page" json:"page"`
	Size int    `form:"size" json:"size"`
}

type PageRequestOld struct { //TODO удалить
	CompanionId string `form:"companionId" json:"companionId"`
	Sort        string `form:"sort" json:"sort"`
	Page        int    `form:"itemPerPag" json:"itemPerPag"`
	Size        int    `form:"offset" json:"offset"`
}

func (p *PageRequest) GetPage() int {
	return p.Page
}
