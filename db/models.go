package db

type User struct {
	StartDate  string
	Email      string
	Name       string
	Password   string
	Session    string
	LastLogin  string
	DeletedAt  string
	GymName    string
	Gender     string
	Permission int
	Salary     int
	Age        int
	ID         int64
}

type LandingPageData struct {
	Title                 string
	StarterSentence       string
	SecondStarterSentence string
	PlansParagraph        string
	AdsOnImageBoldText    string
	AdsOnImageDescription string
	EmailContact          string
	TwitterContact        string
	FacebookContact       string
	InstigramContact      string
	WhatsappContact       string
	ID                    int
}

type LandingPageQNA struct {
	ID int64 `json:"id"`

	LandingPageID int64 `json:"landingPageId"`

	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LandingPageGeneralData struct {
	Title                 string
	StarterSentence       string
	SecondStarterSentence string
	PlansParagraph        string
}

type AdsInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Contacts struct {
	Email     string `json:"email"`
	Twitter   string `json:"twitter"`
	Instagram string `json:"instagram"`
	WhatsApp  string `json:"whatsapp"`
	Facebook  string `json:"facebook"`
}

type Subscriber struct {
	StartedAt     string  `json:"startedAt"`
	Name          string  `json:"name"`
	Surname       string  `json:"surname"`
	DeletedAt     string  `json:"deletedAt" binding:"omitempty"`
	UpdatedAt     string  `json:"updatedAt" binding:"omitempty"`
	CreatedAt     string  `json:"createdAt" binding:"omitempty"`
	EndsAt        string  `json:"endsAt"`
	Gender        string  `json:"gender"`
	Age           int     `json:"age"`
	PaymentAmount float64 `json:"paymentAmount"`
	BucketPrice   float64 `json:"bucketPrice"`
	DaysLeft      int     `json:"daysLeft" binding:"omitempty"`
	Duration      int     `json:"duration" binding:"omitempty"`
	ID            int     `json:"id"`
}

type SubscriberComment struct {
	Sender       *User
	Subscriber   *Subscriber
	Text         string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
	ID           int64
	SenderID     int64
	SubscriberID int64
}

type Message struct {
	Readers *MessageRead
	Text    string
	Sent    string
	ID      int64
}

type MessageRead struct {
	User      *User
	Message   *Message
	ID        int64
	UserID    int64
	MessageID int64
	Read      bool
}

type ProductCategory struct {
	Name     string
	Products []Product
	ID       int64
}

type Product struct {
	Category    *ProductCategory
	Name        string
	Description string
	Marka       string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	ID          int64
	Price       float64
	CategoryID  int64
}

type ProductBasket struct {
	ID       int64
	Quantity int

	CustomerID int64
	Customer   *User

	ProductID int64
	Product   *Product
}

type ExcerciseCategory struct {
	Name       string
	Excercises []Excercise
	ID         int64
}

type Excercise struct {
	Category    *ExcerciseCategory
	Name        string
	Description string
	ID          int64
	CategoryID  int64
}

type Trainer struct {
	Job         string `json:"job"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Instigram   string `json:"instagram"`
	Facebook    string `json:"facebook"`
	Twitter     string `json:"twitter"`
	ID          int64  `json:"id"`
}

type Plan struct {
	Title       string
	Description string
	Duration    string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	Features    []PlanFeature
	ID          int64
	Price       float64
}

type PlanFeature struct {
	Plan   *Plan
	Name   string
	ID     int
	PlanID int
}

type Event struct {
	Actor    *User
	Event    string
	Target   string
	Date     string
	ID       int
	ActorID  int
	TargetID int
}

type SeenEvent struct {
	Event   *Event
	User    *User
	ID      int
	EventID int
	UserID  int
}

type Advice struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Blog struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
	Views       int    `json:"views"`
}
