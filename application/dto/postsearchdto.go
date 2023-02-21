package dto

type PostSearchDto struct {
	AccountIds  []string `form:"accountIds" json:"accountIds"`
	AuthorId    string   `form:"authorId" json:"authorId"`
	Text        string   `form:"text" json:"text"`
	WithFriends bool     `form:"withFriends" json:"withFriends"`
	IsDeleted   bool     `form:"isDeleted" json:"isDeleted"`
	Sort        string   `form:"sort" json:"sort"`
	Page        int      `form:"page" json:"page"`
	Size        int      `form:"size" json:"size"`
}

func (p *PostSearchDto) GetPage() int {
	return p.Page
}
