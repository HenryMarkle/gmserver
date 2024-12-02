package dto

type UpdateLandingPageGeneralInfo_Req struct {
	Title                 string
	StarterSentence       string
	SecondStarterSentence string
	PlansParagraph        string
}

type CreatePlan_Req struct {
	Title       string
	Description string
	Duration    string
	Price       float64
}

type ReplacePlan_Req struct {
	Title       string
	Description string
	Duration    string
	Price       float64
	ID          int64
}

type UpdateAdsInfo_Req struct {
	Title       string
	Description string
}

type CreateProduct_Req struct {
	Name        string
	Description string
	Marka       string
	Price       float64
	CategoryID  int64
}

type UpdateProduct_Req struct {
	Name        string
	Description string
	Marka       string
	Price       float64
	CategoryID  int64
	ID          int64
}

type CreateProductCategory_Req struct {
	Name string
}
