package domain

type LoginResponse struct {
	JwtToken string `json:"token"`
}
