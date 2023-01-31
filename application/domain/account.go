package domain

type Account struct {
	Uid       string   `json:"uid,omitempty"`
	Email     string   `json:"email,omitempty"`
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	Password  string   `json:"password,omitempty"`
	Age       int64    `json:"age,omitempty"`
	DType     []string `json:"dgraph.type,omitempty"`
	//Phone             string    `json:"phone,omitempty"`
	//Photo             string    `json:"photo,omitempty"`
	//PhotoId           string    `json:"photoId,omitempty"`
	//PhotoName         string    `json:"photoName,omitempty"`
	//About             string    `json:"about,omitempty"`
	//City              string    `json:"city,omitempty"`
	//Country           string    `json:"country,omitempty"`
	//StatusCode        string    `json:"statusCode,omitempty"`
	//MessagePermission string    `json:"messagePermission,omitempty"`
	//RegDate           time.Time `json:"regDate,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
	//CreatedOn         time.Time `json:"createdOn,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
	//UpdatedOn         time.Time `json:"updatedOn,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
	//BirthDate         time.Time `json:"birthDate,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
	//LastOnlineTime    time.Time `json:"lastOnlineTime,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
	//IsDeleted         bool      `json:"isDeleted"`
	//IsBlocked         bool      `json:"isBlocked"`
	//IsOnline          bool      `json:"isOnline"`
}

func NewAccount() *Account {
	return &Account{DType: []string{"Account"}}
}
