package domain

type PersonSearchDto struct {
	FirstName string `query:"firstName" pg:"first_name"`
	LastName  string `query:"lastName"  pg:"last_name"`
	//Email        string    `json:"email" pg:"email"`
	//BirthdayTo   time.Time `json:"birthdayTo" pg:"birthday" time_format:"2006-01-02 15:04:05.99Z07:00"`
	//BirthdayFrom time.Time `json:"birthdayFrom" pg:"birthday" time_format:"2006-01-02 15:04:05.99Z07:00"`
}
