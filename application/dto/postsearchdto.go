package dto

type PostSearchDto struct {
	AccountIds  []string `form:"accountIds" json:"accountIds,omitempty"`
	AuthorId    string   `form:"authorId" json:"authorId,omitempty"`
	Text        string   `form:"text" json:"text,omitempty"`
	Author      string   `form:"author" json:"author,omitempty"`
	DateTo      int      `form:"dateTo" json:"dateTo,omitempty"`
	DateFrom    int      `form:"dateFrom" json:"dateFrom,omitempty"`
	Tags        []string `form:"tags" json:"tags,omitempty"`
	WithFriends bool     `form:"withFriends" json:"withFriends,omitempty"`
	Page        int      `form:"page" json:"page,omitempty"`
	Size        int      `form:"size" json:"size,omitempty"`
	IsDeleted   bool     `form:"isDeleted" json:"isDeleted"`
	Sort        string   `form:"sort" json:"sort,omitempty"`
}

func (p *PostSearchDto) GetPage() int {
	return p.Page
}
