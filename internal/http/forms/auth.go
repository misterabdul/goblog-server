package forms

type SignInForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PasswordConfirmForm struct {
	Password string `json:"password" binding:"required"`
}
