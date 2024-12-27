package dto

type CreateComment_Req struct {
	Text         string `json:"text"`
	SenderID     int64  `json:"senderId"`
	SubscriberID int64  `json:"subscriberId"`
}
