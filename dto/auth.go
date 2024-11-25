package dto

type Signin_POST struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword_PATCH struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
