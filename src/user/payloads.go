package user

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
