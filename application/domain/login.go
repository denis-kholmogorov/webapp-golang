package domain

type LoginDto struct {
	Email    string `json:"email" pg:"email"`
	Password string `json:"password" pg:"password"`
}
