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
	account User,
) DBErr {
	dupCheckQuery := `SELECT id FROM User WHERE email = ?;`

	res := db.QueryRow(dupCheckQuery, account.Email)

	var dupId int
	scanErr := res.Scan(&dupId)

	if scanErr != sql.ErrNoRows {
		common.Logger.Printf("Could not add an account with the same email as another account: %v\n", scanErr)
		return ERR_DUPLICATE_ACCOUNT
	}

	createQuery := `INSERT INTO Users (name, email, password, permission, age, gender, salary) VALUES (?, ?, ?, 0, ?, ?)`

	_, execErr := db.Exec(createQuery, account.Name, account.Email, account.Password, account.Age, account.Gender, account.Salary)

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

func GetSubscriberByID(db DB, id int) (*Subscriber, DBErr) {
	query := `SELECT id, name, surname, age, gender, duration, daysLeft, bucketPrice, paymentAmount, startedAt, endsAt, createdAt, updatedAt, deletedAt FROM Subscriber WHERE id = ? WHERE deletedAt IS NULL`

	sub := &Subscriber{}
	scanErr := db.QueryRow(query, id).Scan(&sub.ID, &sub.Name, &sub.Surname, &sub.Age, &sub.Gender, &sub.Duration, &sub.DaysLeft, &sub.BucketPrice, &sub.PaymentAmount, &sub.StartedAt, &sub.EndsAt, &sub.CreatedAt, &sub.UpdatedAt, &sub.DeletedAt)

	if scanErr != nil {
		common.Logger.Printf("Failed to get subscriber ID: %v\n", scanErr)
		return nil, ERR_UNKNOWN
	}

	return sub, ERR_NONE
}

func GetSubscriberByIDWithDeleted(db DB, id int) (*Subscriber, DBErr) {
	query := `SELECT id, name, surname, age, gender, duration, daysLeft, bucketPrice, paymentAmount, startedAt, endsAt, createdAt, updatedAt, deletedAt FROM Subscriber WHERE id = ?`

	sub := &Subscriber{}
	scanErr := db.QueryRow(query, id).Scan(&sub.ID, &sub.Name, &sub.Surname, &sub.Age, &sub.Gender, &sub.Duration, &sub.DaysLeft, &sub.BucketPrice, &sub.PaymentAmount, &sub.StartedAt, &sub.EndsAt, &sub.CreatedAt, &sub.UpdatedAt, &sub.DeletedAt)

	if scanErr != nil {
		common.Logger.Printf("Failed to get subscriber ID: %v\n", scanErr)
		return nil, ERR_UNKNOWN
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

func GetPlans(db DB) ([]Plan, DBErr) {
	query := `SELECT id, title, description, price, duration, createdAt, updatedAt, deletedAt FROM Plan WHERE deletedAt IS NULL`

	rows, execErr := db.Query(query)
	if execErr != nil {
		common.Logger.Printf("Failed to get all plans: %v\n", execErr)
		return nil, ERR_UNKNOWN
	}

	plans := []Plan{}

	counter := 0
	for rows.Next() {
		plan := Plan{}

		scanErr := rows.Scan(&plan.ID, &plan.Title, &plan.Description, &plan.Duration, &plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt)

		if scanErr != nil {
			common.Logger.Printf("Failed to scal plans rows at row (%d): %v\n", counter, scanErr)
			return nil, ERR_UNKNOWN
		}

		plans = append(plans, plan)
	}

	return plans, ERR_NONE
}

func GetPlansWithDeleted(db DB) ([]Plan, DBErr) {
	query := `SELECT id, title, description, price, duration, createdAt, updatedAt, deletedAt FROM Plan`

	rows, execErr := db.Query(query)
	if execErr != nil {
		common.Logger.Printf("Failed to get all plans: %v\n", execErr)
		return nil, ERR_UNKNOWN
	}

	plans := []Plan{}

	counter := 0
	for rows.Next() {
		plan := Plan{}

		scanErr := rows.Scan(&plan.ID, &plan.Title, &plan.Description, &plan.Duration, &plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt)

		if scanErr != nil {
			common.Logger.Printf("Failed to scal plans rows at row (%d): %v\n", counter, scanErr)
			return nil, ERR_UNKNOWN
		}

		plans = append(plans, plan)
	}

	return plans, ERR_NONE
}

func GetPlanByID(db DB, id int) (*Plan, DBErr) {
	query := `SELECT title, description, price, duration, createdAt, updatedAt, deletedAt FROM Plan WHERE id = ?`

	plan := &Plan{}

	scanErr := db.QueryRow(query, id).Scan(&plan.Title, &plan.Description, &plan.Price, &plan.Duration, &plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt)

	if scanErr != nil {
		common.Logger.Printf("Failed to get a plan by ID: %v\n", scanErr)
		return nil, ERR_UNKNOWN
	}

	return plan, ERR_NONE
}

func CreatePlan(db DB, plan Plan) (int, DBErr) {
	query := `
  BEGIN TRANSACTION;

  INSERT INTO Plan (title, description, price, duration) VALUES (?, ?, ?, ?);

  SELECT LAST_INSERT_ROWID();

  COMMIT;
  `

	var id int
	scanErr := db.QueryRow(query, plan.Title, plan.Description, plan.Price, plan.Duration).Scan(&id)

	if scanErr != nil {
		common.Logger.Printf("Failed to create a plan: %v\n", scanErr)
		return 0, ERR_UNKNOWN
	}

	return id, ERR_NONE
}

func DeletePlanByID(db DB, id int) DBErr {
	query := `DELETE FROM Plan WHERE id = ?`

	_, err := db.Exec(query, id)
	if err != nil {
		common.Logger.Printf("Failed to delete a plan by ID: %v\n", err)
		return ERR_UNKNOWN
	}

	return ERR_NONE
}

func ReplacePlan(db DB, plan Plan) DBErr {
	tx, txErr := db.Begin()

	if tx != nil {
		common.Logger.Printf("Could not begin a new transaction to replace a plan: %v\n", txErr)
		return ERR_UNKNOWN
	}

	deleteQuery := `DELETE FROM Plan WHERE id = ?`

	_, deleteErr := tx.Exec(deleteQuery, plan.ID)

	if deleteErr != nil {
		tx.Rollback()
		common.Logger.Printf("Failed to delete a plan by ID in a transaction: %v\n", deleteErr)
		return ERR_UNKNOWN
	}

	insertQuery := `INSERT INTO Plan (title, description, price, duration) VALUES (?, ?, ?, ?);`

	_, insertErr := tx.Exec(insertQuery, plan.Title, plan.Description, plan.Price, plan.Duration)

	if insertErr != nil {
		tx.Rollback()
		common.Logger.Printf("Failed to create a plan in a transaction: %v\n", insertErr)
		return ERR_UNKNOWN
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		tx.Rollback()
		common.Logger.Printf("Could not commit plan replacement: %v\n", commitErr)
		return ERR_UNKNOWN
	}

	return ERR_NONE
}
