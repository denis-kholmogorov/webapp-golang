package domain

import "time"

type Account struct {
	tableName struct{}  `pg:"public.account"`
	Id        int64     `json:"id" pg:"id"`
	Age       int64     `json:"age" pg:"age"`
	FirstName string    `json:"firstName" pg:"first_name"`
	LastName  string    `json:"lastName" pg:"last_name"`
	Email     string    `json:"email" pg:"email"`
	Password  string    `json:"password" pg:"password"`
	Birthday  time.Time `json:"birthday" pg:"birthday" time_format:"2006-01-02 15:04:05.99Z07:00"`
}
