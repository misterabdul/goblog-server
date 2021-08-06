package forms

type UpdateMeForm struct {
	FirstName string `form:"firstname" json:"firstname" xml:"firstname"`
	LastName  string `form:"lastname" json:"lastname" xml:"lastname"`
	Username  string `form:"username" json:"username" xml:"username"`
	Email     string `form:"email" json:"email" xml:"email" binding:"email"`
}

type UpdateMePasswordForm struct {
	NewPassword        string `form:"newpassword" json:"newpassword" xml:"newpassword" bind:"required"`
	NewPasswordConfirm string `form:"newpasswordconfirm" json:"newpasswordconfirm" xml:"newpasswordconfirm" bind:"required"`
}
