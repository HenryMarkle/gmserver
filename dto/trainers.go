package dto

type CreateTrainer_Req struct {
	Job         string `json:"job"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Instigram   string `json:"instigram"`
	Facebook    string `json:"facebook"`
	Twitter     string `json:"twitter"`
}

type UpdateTrainer_Req struct {
	Job         string `json:"job"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Instigram   string `json:"instigram"`
	Facebook    string `json:"facebook"`
	Twitter     string `json:"twitter"`
	ID          int64  `json:"id"`
}
