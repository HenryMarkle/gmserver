package dto

type CreateBlog_Req struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Views       int    `json:"views,omitempty"`
}

type UpdateBlog_Req struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Views       int    `json:"views,omitempty"`
}

type GetBlog_Res struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Views       int    `json:"views"`
}
