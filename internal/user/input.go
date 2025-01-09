package user

type LoginInputUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordInputUser struct {
	Password string `json:"password" binding:"required"`
}
