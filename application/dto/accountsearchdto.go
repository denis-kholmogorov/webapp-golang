package dto

type AccountSearchDto struct {
	Id           string   `form:"id" json:"id"`
	Ids          []string `form:"ids" json:"ids"`
	BlockedByIds []string `form:"blockedByIds" json:"blockedByIds"`
	Author       string   `form:"author" json:"author"`
	FirstName    string   `form:"firstName" json:"firstName"`
	LastName     string   `form:"lastName" json:"lastName"`
	City         string   `form:"city" json:"city"`
	Country      string   `form:"country" json:"country"`
	AgeTo        int      `form:"ageTo" json:"ageTo"`
	AgeFrom      int      `form:"ageFrom" json:"ageFrom"`
	IsDeleted    bool     `form:"isDeleted" json:"isDeleted"`
	IsBlocked    bool     `form:"isBlocked" json:"isBlocked"`
	Page         int      `form:"page" json:"page"`
	Size         int      `form:"size" json:"size"`
}

func (p *AccountSearchDto) GetPage() int {
	return p.Page
}
