package user

type LoginInputUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdatePasswordInputUser struct {
	Password string `json:"password" validate:"required"`
}
