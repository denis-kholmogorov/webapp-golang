package domain

import "time"

type Person struct {
	tableName struct{}  `pg:"public.person"`
	Id        int       `json:"id" pg:"id"`
	Age       int       `json:"age" pg:"age"`
	FirstName string    `json:"firstName" pg:"first_name"`
	Birthday  time.Time `json:"birthday" pg:"birthday" time_format:"2006-01-02 15:04:05.99Z07:00"`
}
