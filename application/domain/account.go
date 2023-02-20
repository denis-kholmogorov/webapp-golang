package domain

import "time"

type Account struct {
	Uid               string     `json:"id,omitempty"`
	Email             string     `json:"email,omitempty"`
	FirstName         string     `json:"firstName,omitempty"`
	LastName          string     `json:"lastName,omitempty"`
	Password          string     `json:"password,omitempty"`
	Age               int64      `json:"age,omitempty"`
	IsDeleted         bool       `json:"isDeleted"`
	IsBlocked         bool       `json:"isBlocked"`
	IsOnline          bool       `json:"isOnline"`
	DType             []string   `json:"dgraph.type,omitempty"`
	Posts             []Post     `json:"posts,omitempty"`
	Phone             string     `json:"phone,omitempty"`
	Photo             string     `json:"photo,omitempty"`
	PhotoId           string     `json:"photoId,omitempty"`
	PhotoName         string     `json:"photoName,omitempty"`
	About             string     `json:"about,omitempty"`
	City              string     `json:"city,omitempty"`
	Country           string     `json:"country,omitempty"`
	StatusCode        string     `json:"statusCode,omitempty"`
	MessagePermission string     `json:"messagePermission,omitempty"`
	CreatedOn         *time.Time `json:"createdOn" time_format:"2006-01-02 15:04:05.99Z07:00"`
	UpdatedOn         *time.Time `json:"updatedOn" time_format:"2006-01-02 15:04:05.99Z07:00"`
	BirthDate         *time.Time `json:"birthDate,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
	LastOnlineTime    *time.Time `json:"lastOnlineTime,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
}

type PageResponse struct {
	Content      []any   `json:"content,omitempty"`
	TotalElement int     `json:"totalElement,omitempty"`
	TotalPages   int     `json:"totalPages,omitempty"`
	Number       int     `json:"number"`
	Size         int     `json:"size,omitempty"`
	Count        []Count `json:"count,omitempty"`
}

type Count struct {
	TotalElement int `json:"totalElement,omitempty"`
}
