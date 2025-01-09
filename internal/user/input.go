package user

type LoginInputUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdatePasswordInputUser struct {
	Password string `json:"password"`
}
