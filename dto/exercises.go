package dto

type ExcerciseCategory_Res struct {
	Name      string          `json:"name"`
	Exercises []Excercise_Res `json:"exercises"`
	ID        int64           `json:"id"`
}

type Excercise_Res struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"categoryID"`
}

type CreateExcersizeCategory_Req struct {
	Name string `json:"name"`
}

type CreateExcercise_Req struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  int64  `json:"categoryId"`
}

type UpdateExcercise_Req struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"categoryID"`
}
