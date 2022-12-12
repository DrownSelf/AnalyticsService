package handlers

type LogInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
