package dto

type CreateAdvice_Req struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateAdvice_Req struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
