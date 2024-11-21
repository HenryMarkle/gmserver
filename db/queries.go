package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/HenryMarkle/gmserver/common"
	_ "github.com/mattn/go-sqlite3"
)

type (
	DB    = *sql.DB
	DBErr = int
)

const (
	ERR_NONE DBErr = iota
	ERR_DUPLICATE_ACCOUNT

	ERR_UNKNOWN
)

func Construct(db DB) error {
	initPath := "init.sql"

	initScript, readErr := os.ReadFile(initPath)
	if readErr != nil {
		common.Logger.Printf("Failed to read SQLite init file: %v", readErr)
		return readErr
	}

	_, execErr := db.Exec(string(initScript))

	if execErr != nil {
		common.Logger.Printf("Failed to execute SQLite init script: %v", execErr)
		return execErr
	}

	return nil
}

func AddAccount(
	db DB,
	name, email, password string,
	age, salary, gender int,
) DBErr {
	dupCheckQuery := `SELECT id FROM User WHERE email = ?;`

	res := db.QueryRow(dupCheckQuery, email)

	var dupId int
	scanErr := res.Scan(&dupId)

	if scanErr != sql.ErrNoRows {
		common.Logger.Printf("Could not add an account with the same email as another account: %v\n", scanErr)
		return ERR_DUPLICATE_ACCOUNT
	}

	createQuery := `INSERT INTO Users (name, email, password, permission, age, gender, salary) VALUES (?, ?, ?, 0, ?, ?)`

	_, execErr := db.Exec(createQuery, name, email, password, age, gender, salary)

	if execErr != nil {
		common.Logger.Printf("Failed to create account: %v\n", execErr)
		return ERR_UNKNOWN
	}

	return ERR_NONE
}

func GetTotalSubscriberPaymentAmount(db DB) (int, DBErr) {
	query := `SELECT SUM(paymentAmount) FROM Subscriber`

	var totalAmount int

	err := db.QueryRow(query).Scan(&totalAmount)
	if err != nil {
		common.Logger.Printf("Failed to query the total paid amount by subscribers: %v\n", err)
		return 0, ERR_UNKNOWN
	}

	return totalAmount, ERR_NONE
}

func GetSubscriberCount(db DB) (int, DBErr) {
	query := `SELECT COUNT(*) FROM Subscriber`

	var number int

	err := db.QueryRow(query).Scan(&number)
	if err != nil {
		common.Logger.Printf("Failed to count the total number of subscribers: %v\n", err)
		return 0, ERR_UNKNOWN
	}

	return number, ERR_NONE
}

// / Time string must be of format '2024-11-21 12:00:00'
func GetAllSubscribersEndingBefore(db DB, time string) (int, DBErr) {
	query := `SELECT COUNT(endsAt) FROM Subscriber WHERE endsAt < ?`

	var number int

	err := db.QueryRow(query, time).Scan(&number)
	if err != nil {
		common.Logger.Printf("Failed to count subscribers ending before a given fate: %v\n", err)
		return 0, ERR_UNKNOWN
	}

	return number, ERR_NONE
}

// / Time string must be of format '2024-11-21 12:00:00'
func GetAllExpiredSubscribers(db DB, time string) (int, DBErr) {
	query := `SELECT COUNT(endsAt) FROM Subscriber WHERE endsAt > ?`

	var number int

	err := db.QueryRow(query, time).Scan(&number)
	if err != nil {
		common.Logger.Printf("Failed to count expired subscibers: %v\n", err)
		return 0, ERR_UNKNOWN
	}

	return number, ERR_NONE
}

func CreateSubscriber(db DB,
	name, surname, startDate, endDate string,
	age, gender, payment, bucketPrice int,
) DBErr {
	query := `
  INSERT INTO Subscriber 
  (name, surname, age, paymentAmount, startedAt, endsAt, bucketPrice) 
  VALUES 
  (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, name, surname, age, payment, startDate, endDate, bucketPrice)
	if err != nil {
		common.Logger.Printf("Failed to create subscriber: %v\n", err)
		return ERR_UNKNOWN
	}

	return ERR_NONE
}

func GetAllSubscribers(db DB, limit int) ([]Subscriber, DBErr) {
	query := `SELECT id, name, surname, age, gender, duration, daysLeft, bucketPrice, paymentAmount, startedAt, endsAt, createdAt, updatedAt, deletedAt FROM Subscriber`

	if limit > 0 {
		query = fmt.Sprintf("%s %s %d", query, "LIMIT", limit)
	}

	rows, err := db.Query(query)
	if err != nil {
		common.Logger.Printf("Failed to get all subscribers: %v\n", err)
		return nil, ERR_UNKNOWN
	}

	defer rows.Close()

	subs := []Subscriber{}

	counter := 0

	for rows.Next() {
		var (
			id, age, duration, daysLeft int
			name, surname, gender,
			startedAt, endsAt, createdAt, updatedAt, deletedAt string

			bucketPrice, paymentAmount float64
		)

		scanErr := rows.Scan(&id, &name, &surname, &age, &gender, &duration, &daysLeft, &bucketPrice, &paymentAmount, &startedAt, &endsAt, &createdAt, &deletedAt)

		if scanErr != nil {
			common.Logger.Printf("Failed to scan subscriber rows at row (%d): %v\n", counter, scanErr)
		}

		subs = append(subs, Subscriber{
			ID:            id,
			Name:          name,
			Surname:       surname,
			Age:           age,
			Gender:        gender,
			Duration:      duration,
			DaysLeft:      daysLeft,
			BucketPrice:   bucketPrice,
			PaymentAmount: paymentAmount,
			StartedAt:     startedAt,
			EndsAt:        endsAt,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			DeletedAt:     deletedAt,
		})

		counter++
	}

	return subs, ERR_NONE
}

func GetSubscriberByID(db DB, id int, includeDeleted bool) (*Subscriber, DBErr) {
	query := `SELECT id, name, surname, age, gender, duration, daysLeft, bucketPrice, paymentAmount, startedAt, endsAt, createdAt, updatedAt, deletedAt FROM Subscriber WHERE id = ?`

	if !includeDeleted {
		query = fmt.Sprintf("%s AND deletedAt IS NULL", query)
	}

	var (
		_id, age, duration, daysLeft int
		name, surname, gender,
		startedAt, endsAt, createdAt, updatedAt, deletedAt string

		bucketPrice, paymentAmount float64
	)

	scanErr := db.QueryRow(query, id).Scan(&_id, &name, &surname, &age, &gender, &duration, &daysLeft, &bucketPrice, &paymentAmount, &startedAt, &endsAt, &createdAt, &updatedAt, &deletedAt)

	if scanErr != nil {
		common.Logger.Printf("Failed to get subscriber ID: %v\n", scanErr)
		return nil, ERR_UNKNOWN
	}

	sub := &Subscriber{
		ID:            id,
		Name:          name,
		Surname:       surname,
		Age:           age,
		Gender:        gender,
		Duration:      duration,
		DaysLeft:      daysLeft,
		BucketPrice:   bucketPrice,
		PaymentAmount: paymentAmount,
		StartedAt:     startedAt,
		EndsAt:        endsAt,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		DeletedAt:     deletedAt,
	}

	return sub, ERR_NONE
}

func DeleteSubscriberByID(db DB, id int, permanent bool) DBErr {
	query := `DELETE FROM Subscriber WHERE id = ?`

	if !permanent {
		query = `UPDATE Subscriber SET deletedAt = CURRENT_TIMESTAMP WHERE id = ?`
	}

	_, execErr := db.Exec(query, id)

	if execErr != nil {
		common.Logger.Printf("Failed to delete a subscriber by ID: %v\n", execErr)
		return ERR_UNKNOWN
	}

	return ERR_NONE
}

func UpdateSubscriber(db DB, data Subscriber) DBErr {
	query := `
  UPDATE Subscriber 
  SET 
    name = ?,
    surname = ?,
    age = ?,
    gender = ?,
    duration = ?,
    daysLeft = ?,
    bucketPrice = ?,
    paymentAmount = ?,
    startedAt = ?,
    endsAt = ?,
    createdAt = ?,
    updatedAt = ?,
    deletedAt = ?
  WHERE id = ?`

	_, execErr := db.Exec(query,
		data.Name,
		data.Surname,
		data.Age,
		data.Gender,
		data.Duration,
		data.DaysLeft,
		data.BucketPrice,
		data.PaymentAmount,
		data.StartedAt,
		data.EndsAt,
		data.CreatedAt,
		data.UpdatedAt,
		data.DeletedAt,
		data.ID)

	if execErr != nil {
		common.Logger.Printf("Failed to update a subscriber: %v\n", execErr)
		return ERR_UNKNOWN
	}

	return ERR_NONE
}
