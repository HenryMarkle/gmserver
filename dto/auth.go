package dto

type Signin_Req struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword_Req struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type CreateUser_Req struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	StartDate string `json:"startDate"`
	//	Permission int    `json:"permission"`
	Age    int `json:"age"`
	Salary int `json:"salary"`
}
