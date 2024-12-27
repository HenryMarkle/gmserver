package dto

type CreateSubscriber_Req struct {
	Name          string  `json:"name"`
	Surname       string  `json:"surname"`
	StartedAt     string  `json:"startedAt"`
	DeletedAt     string  `json:"deletedAt"`
	UpdatedAt     string  `json:"updatedAt"`
	EndsAt        string  `json:"endsAt"`
	Gender        string  `json:"gender"`
	Age           int     `json:"age"`
	PaymentAmount float64 `json:"paymentAmount"`
	BucketPrice   float64 `json:"bucketPrice"`
	DaysLeft      int     `json:"daysLeft"`
	Duration      int     `json:"duration"`
}
