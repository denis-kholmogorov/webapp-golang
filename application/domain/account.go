package domain

import (
	"time"
)

type Account struct {
	tableName         struct{}  `pg:"public.account"`
	Id                int64     `json:"id" pg:"id"`
	Email             string    `json:"email" pg:"email"`
	FirstName         string    `json:"firstName" pg:"first_name"`
	LastName          string    `json:"lastName" pg:"last_name"`
	Password          string    `json:"password" pg:"password"`
	Age               int64     `json:"age" pg:"age"`
	Phone             string    `json:"phone" pg:"phone"`
	Photo             string    `json:"photo" pg:"photo"`
	PhotoId           string    `json:"photoId" pg:"photo_id"`
	PhotoName         string    `json:"photoName" pg:"photo_name"`
	About             string    `json:"about" pg:"about"`
	City              string    `json:"city" pg:"city"`
	Country           string    `json:"country" pg:"country"`
	StatusCode        string    `json:"statusCode" pg:"status_code"`
	MessagePermission string    `json:"messagePermission" pg:"message_permission"`
	RegDate           time.Time `json:"regDate" pg:"reg_date" time_format:"2006-01-02 15:04:05.99Z07:00"`
	CreatedOn         time.Time `json:"createdOn" pg:"created" time_format:"2006-01-02 15:04:05.99Z07:00"`
	UpdatedOn         time.Time `json:"updatedOn" pg:"updated" time_format:"2006-01-02 15:04:05.99Z07:00"`
	BirthDate         time.Time `json:"birthDate" pg:"birth_date" time_format:"2006-01-02 15:04:05.99Z07:00"`
	LastOnlineTime    time.Time `json:"lastOnlineTime" pg:"last_online_time" time_format:"2006-01-02 15:04:05.99Z07:00"`
	IsDeleted         bool      `json:"isDeleted" pg:"is_deleted"`
	IsBlocked         bool      `json:"isBlocked" pg:"is_blocked"`
	IsOnline          bool      `json:"isOnline" pg:"is_online"`
}
