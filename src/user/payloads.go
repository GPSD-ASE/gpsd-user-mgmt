package user

type signInRequest struct {
	Username string `json"user_name"`
	Password string `json"pass"`
}
