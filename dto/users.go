package dto

type UpdateUser_Req struct {
	StartDate string `json:"startDate"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	GymName   string `json:"gymName"`
	Gender    string `json:"gender"`
	Salary    int    `json:"salary"`
	Age       int    `json:"age"`
	ID        int64  `json:"id"`
}

type User_Res struct {
	StartDate  string `json:"startDate"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GymName    string `json:"gymName"`
	Gender     string `json:"gender"`
	LastLogin  string `json:"lastLogin"`
	Permission int    `json:"permission"`
	Salary     int    `json:"salary"`
	Age        int    `json:"age"`
	ID         int64  `json:"id"`
}
