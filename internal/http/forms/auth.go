package forms

type SignInForm struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type PasswordConfirmForm struct {
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
