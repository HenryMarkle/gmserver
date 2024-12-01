package dto

type CreateSubscriber_Req struct {
	Name          string
	Surname       string
	StartedAt     string
	DeletedAt     string
	UpdatedAt     string
	EndsAt        string
	Gender        string
	Age           int
	PaymentAmount float64
	BucketPrice   float64
	DaysLeft      int
	Duration      int
}
