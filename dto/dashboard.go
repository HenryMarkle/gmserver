package dto

type UpdateLandingPageGeneralInfo_Req struct {
	Title                 string `json:"title"`
	StarterSentence       string `json:"starterSentence"`
	SecondStarterSentence string `json:"secondStarterSentence"`
	PlansParagraph        string `json:"plansParagraph"`
}

type CreatePlan_Req struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    string  `json:"duration"`
	Price       float64 `json:"price"`
}

type ReplacePlan_Req struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    string  `json:"duration"`
	Price       float64 `json:"price"`
	ID          int64   `json:"id"`
}

type UpdateAdsInfo_Req struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateProduct_Req struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Marka       string  `json:"marka"`
	Price       float64 `json:"price"`
	CategoryID  int64   `json:"categoryId"`
}

type UpdateProduct_Req struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Marka       string  `json:"marka"`
	Price       float64 `json:"price"`
	CategoryID  int64   `json:"categoryId"`
	ID          int64   `json:"id"`
}

type CreateProductCategory_Req struct {
	Name string `json:"name"`
}

type CreateQNA_Req struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
