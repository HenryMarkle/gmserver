package dto

type UpdateBlog_Req struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
	Views       int    `json:"views"`
}
