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

type LandingPageGeneralData struct {
	Title                 string
	StarterSentence       string
	SecondStarterSentence string
	PlansParagraph        string
}

type AdsInfo struct {
	Title       string
	Description string
}

type Contacts struct {
	Email     string
	Twitter   string
	Instagram string
	WhatsApp  string
	Facebook  string
}

type Subscriber struct {
	StartedAt     string
	Name          string
	Surname       string
	DeletedAt     string
	UpdatedAt     string
	CreatedAt     string
	EndsAt        string
	Gender        string
	Age           int
	PaymentAmount float64
	BucketPrice   float64
	DaysLeft      int
	Duration      int
	ID            int
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
	ID       int
}

type Product struct {
	Category    *ProductCategory
	Name        string
	Description string
	Marka       string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	ID          int
	Price       float64
	CategoryID  int
}

type ExcerciseCategory struct {
	Name       string
	Excercises []Excercise
	ID         int
}

type Excercise struct {
	Category    *ExcerciseCategory
	Name        string
	Description string
	ID          int
	CategoryID  int
}

type Trainer struct {
	Job         string
	Name        string
	Description string
	Instigram   string
	Facebook    string
	Twitter     string
	ID          int64
}

type Plan struct {
	Title       string
	Description string
	Duration    string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	ID          int
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
