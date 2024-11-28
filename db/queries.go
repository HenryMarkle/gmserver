package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/HenryMarkle/gmserver/common"
	_ "github.com/go-sql-driver/mysql"
)

func Construct(db *sql.DB) error {
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
	db *sql.DB,
	account User,
) error {
	dupCheckQuery := `SELECT id FROM User WHERE email = ?;`

	res := db.QueryRow(dupCheckQuery, account.Email)

	var dupId int
	scanErr := res.Scan(&dupId)

	if scanErr != sql.ErrNoRows {
		return fmt.Errorf("Could not add an account with the same email as another account: %w\n", scanErr)
	}

	createQuery := `INSERT INTO Users (name, email, password, permission, age, gender, salary) VALUES (?, ?, ?, 0, ?, ?)`

	_, execErr := db.Exec(createQuery, account.Name, account.Email, account.Password, account.Age, account.Gender, account.Salary)

	if execErr != nil {
		return fmt.Errorf("Failed to create account: %w\n", execErr)
	}

	return nil
}

func GetUserByID(db *sql.DB, id int64) (*User, error) {
	query := `SELECT email, name, password, session, lastLogin, age, salary, permission, gender, startDate FROM User WHERE id = ?`

	row := db.QueryRow(query, id)

	user := User{ID: id}

	scanErr := row.Scan(&user.Email, &user.Name, &user.Password, &user.Session, &user.LastLogin, &user.Age, &user.Salary, &user.Permission, &user.Gender, &user.StartDate)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("Failed to get a user by ID (id: %d): %w\n", id, scanErr)
	}

	return &user, nil
}

func GetUserBySession(db *sql.DB, session string) (*User, error) {
	query := `SELECT id, email, name, password, lastLogin, age, salary, permission, gender, startDate FROM User WHERE session = ?`

	row := db.QueryRow(query, session)

	user := User{Session: session}

	scanErr := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Session, &user.LastLogin, &user.Age, &user.Salary, &user.Permission, &user.Gender, &user.StartDate)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("Failed to get a user by session: %w\n", scanErr)
	}

	return &user, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `SELECT id, session, name, password, lastLogin, age, salary, permission, gender, startDate FROM User WHERE email = ?`

	row := db.QueryRow(query, email)

	user := User{Email: email}

	scanErr := row.Scan(&user.ID, &user.Session, &user.Name, &user.Password, &user.Session, &user.LastLogin, &user.Age, &user.Salary, &user.Permission, &user.Gender, &user.StartDate)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("Failed to get a user by session: %w\n", scanErr)
	}

	return &user, nil
}

func UserExistsByID(db *sql.DB, id int64) bool {
	query := `SELECT EXISTS (
    SELECT 1 FROM User WHERE id = ?
  );`

	exists := false

	_ = db.QueryRow(query, id).Scan(&exists)

	return exists
}

func UserExistsByEmail(db *sql.DB, email string) bool {
	query := `SELECT EXISTS (
    SELECT 1 FROM User WHERE email = ?
  );`

	exists := false

	_ = db.QueryRow(query, email).Scan(&exists)

	return exists
}

func UpdateUser(db *sql.DB, data User) error {
	query := `UPDATE User SET email = ?, name = ?, gymName = ?, age = ?, startDate = ?, salary = ?, gender = ? WHERE id = ?`

	_, execErr := db.Exec(query, data.Email, data.Name, data.GymName, data.Age, data.StartDate, data.Salary, data.Gender, data.ID)
	if execErr != nil {
		return fmt.Errorf("Failed to update user by ID (id: %d): %w\n", data.ID, execErr)
	}

	return nil
}

func DeleteUserByID(db *sql.DB, id int64) error {
	query := `DELETE FROM User WHERE id = ?`
	_, execErr := db.Exec(query, id)
	if execErr != nil {
		return fmt.Errorf("Failed to delete a user by ID (id: %d): %w\n", id, execErr)
	}

	return nil
}

func CountUsers(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM User`
	var count int
	scanErr := db.QueryRow(query).Scan(&count)
	if scanErr != nil {
		return 0, fmt.Errorf("Failed to count users: %w\n", scanErr)
	}
	return count, nil
}

func MarkUserAsDeleted(db *sql.DB, id int64) error {
	query := `UPDATE User SET deletedAt = CURRENT_TIMESTAMP WHERE id = ?`
	_, execErr := db.Exec(query, id)
	if execErr != nil {
		return fmt.Errorf("Failed to mark a user by ID (id: %d): %w\n ", id, execErr)
	}
	return nil
}

func ChangeUserPassword(db *sql.DB, id int64, newPassword string) error {
	query := `UPDATE User SET password = ? WHERE id = ?`

	_, execErr := db.Exec(query, newPassword, id)
	if execErr != nil {
		return fmt.Errorf("Failed to update user password by ID (id: %d): %w\n", id, execErr)
	}

	return nil
}

func ChangeGymName(db *sql.DB, id int64, newGymName string) error {
	query := `UPDATE User SET gymName = ? WHERE id = ?`
	_, execErr := db.Exec(query, newGymName, id)
	if execErr != nil {
		return fmt.Errorf("Failed to update gym name of a user (id: %d): %w\n", id, execErr)
	}

	return nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, email, name, gender, age, salary, startDate, permission FROM User WHERE deletedAt IS NULL`

	rows, queryErr := db.Query(query)
	if queryErr != nil {
		return nil, fmt.Errorf("Failed to get all users: %w\n", queryErr)
	}

	defer rows.Close()

	users := []User{}
	counter := 0

	for rows.Next() {
		user := User{}
		scanErr := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Gender, &user.Age, &user.Salary, &user.StartDate, &user.Permission)

		if scanErr != nil {
			common.Logger.Printf("Failed to scan a user from rows at row (%d): %v\n", counter, scanErr)
		} else {
			users = append(users, user)
		}

		counter++
	}

	return users, nil
}

func GetTotalSalaries(db *sql.DB) (int, error) {
	query := `SELECT SUM(salary) FROM User`

	sum := 0
	scanErr := db.QueryRow(query).Scan(&sum)
	if scanErr != nil {
		return 0, fmt.Errorf("Failed to sum the total salaries: %w\n", scanErr)
	}

	return sum, nil
}

func GetTotalSubscriberPaymentAmount(db *sql.DB) (int, error) {
	query := `SELECT SUM(paymentAmount) FROM Subscriber`

	var totalAmount int

	err := db.QueryRow(query).Scan(&totalAmount)
	if err != nil {
		return 0, fmt.Errorf("Failed to query the total paid amount by subscribers: %w\n", err)
	}

	return totalAmount, nil
}

func GetSubscriberCount(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM Subscriber`

	var number int

	err := db.QueryRow(query).Scan(&number)
	if err != nil {
		return 0, fmt.Errorf("Failed to count the total number of subscribers: %w\n", err)
	}

	return number, nil
}

// / Time string must be of format '2024-11-21 12:00:00'
func GetAllSubscribersEndingBefore(db *sql.DB, time string) (int, error) {
	query := `SELECT COUNT(endsAt) FROM Subscriber WHERE endsAt < ?`

	var number int

	err := db.QueryRow(query, time).Scan(&number)
	if err != nil {
		return 0, fmt.Errorf("Failed to count subscribers ending before a given fate: %w\n", err)
	}

	return number, nil
}

// / Time string must be of format '2024-11-21 12:00:00'
func GetAllExpiredSubscribers(db *sql.DB, time string) (int, error) {
	query := `SELECT COUNT(endsAt) FROM Subscriber WHERE endsAt > ?`

	var number int

	err := db.QueryRow(query, time).Scan(&number)
	if err != nil {
		return 0, fmt.Errorf("Failed to count expired subscibers: %w\n", err)
	}

	return number, nil
}

func CreateSubscriber(db *sql.DB,
	name, surname, startDate, endDate string,
	age, gender, payment, bucketPrice int,
) error {
	query := `
  INSERT INTO Subscriber 
  (name, surname, age, paymentAmount, startedAt, endsAt, bucketPrice) 
  VALUES 
  (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, name, surname, age, payment, startDate, endDate, bucketPrice)
	if err != nil {
		return fmt.Errorf("Failed to create subscriber: %w\n", err)
	}

	return nil
}

func GetAllSubscribers(db *sql.DB, limit int) ([]Subscriber, error) {
	query := `SELECT id, name, surname, age, gender, duration, daysLeft, bucketPrice, paymentAmount, startedAt, endsAt, createdAt, updatedAt, deletedAt FROM Subscriber`

	if limit > 0 {
		query = fmt.Sprintf("%s %s %d", query, "LIMIT", limit)
	}

	rows, err := db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Subscriber{}, nil
		}

		return nil, fmt.Errorf("Failed to get all subscribers: %w\n", err)
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
		} else {
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
		}

		counter++
	}

	return subs, nil
}

func GetSubscriberByID(db *sql.DB, id int64) (*Subscriber, error) {
	query := `SELECT id, name, surname, age, gender, duration, daysLeft, bucketPrice, paymentAmount, startedAt, endsAt, createdAt, updatedAt, deletedAt FROM Subscriber WHERE id = ? WHERE deletedAt IS NULL`

	sub := &Subscriber{}
	scanErr := db.QueryRow(query, id).Scan(&sub.ID, &sub.Name, &sub.Surname, &sub.Age, &sub.Gender, &sub.Duration, &sub.DaysLeft, &sub.BucketPrice, &sub.PaymentAmount, &sub.StartedAt, &sub.EndsAt, &sub.CreatedAt, &sub.UpdatedAt, &sub.DeletedAt)

	if scanErr != nil {
		return nil, fmt.Errorf("Failed to get subscriber ID: %w\n", scanErr)
	}

	return sub, nil
}

func GetSubscriberByIDWithDeleted(db *sql.DB, id int) (*Subscriber, error) {
	query := `SELECT id, name, surname, age, gender, duration, daysLeft, bucketPrice, paymentAmount, startedAt, endsAt, createdAt, updatedAt, deletedAt FROM Subscriber WHERE id = ?`

	sub := &Subscriber{}
	scanErr := db.QueryRow(query, id).Scan(&sub.ID, &sub.Name, &sub.Surname, &sub.Age, &sub.Gender, &sub.Duration, &sub.DaysLeft, &sub.BucketPrice, &sub.PaymentAmount, &sub.StartedAt, &sub.EndsAt, &sub.CreatedAt, &sub.UpdatedAt, &sub.DeletedAt)

	if scanErr != nil {
		return nil, fmt.Errorf("Failed to get subscriber ID: %w\n", scanErr)
	}

	return sub, nil
}

func DeleteSubscriberByID(db *sql.DB, id int, permanent bool) error {
	query := `DELETE FROM Subscriber WHERE id = ?`

	if !permanent {
		query = `UPDATE Subscriber SET deletedAt = CURRENT_TIMESTAMP WHERE id = ?`
	}

	_, execErr := db.Exec(query, id)

	if execErr != nil {
		return fmt.Errorf("Failed to delete a subscriber by ID: %w\n", execErr)
	}

	return nil
}

func UpdateSubscriber(db *sql.DB, data Subscriber) error {
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
		return fmt.Errorf("Failed to update a subscriber: %w\n", execErr)
	}

	return nil
}

func GetPlans(db *sql.DB) ([]Plan, error) {
	query := `SELECT id, title, description, price, duration, createdAt, updatedAt, deletedAt FROM Plan WHERE deletedAt IS NULL`

	rows, execErr := db.Query(query)
	if execErr != nil {
		if execErr == sql.ErrNoRows {
			return []Plan{}, nil
		}

		return nil, fmt.Errorf("Failed to get all plans: %w\n", execErr)
	}

	defer rows.Close()

	plans := []Plan{}

	counter := 0
	for rows.Next() {
		plan := Plan{}

		scanErr := rows.Scan(&plan.ID, &plan.Title, &plan.Description, &plan.Duration, &plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt)

		if scanErr != nil {
			common.Logger.Printf("Failed to scal plans rows at row (%d): %v\n", counter, scanErr)
		} else {
			plans = append(plans, plan)
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

func GetPlansWithDeleted(db *sql.DB) ([]Plan, error) {
	query := `SELECT id, title, description, price, duration, createdAt, updatedAt, deletedAt FROM Plan`

	rows, execErr := db.Query(query)
	if execErr != nil {
		if execErr == sql.ErrNoRows {
			return []Plan{}, nil
		}

		return nil, fmt.Errorf("Failed to get all plans: %w\n", execErr)
	}

	defer rows.Close()

	plans := []Plan{}

	counter := 0
	for rows.Next() {
		plan := Plan{}

		scanErr := rows.Scan(&plan.ID, &plan.Title, &plan.Description, &plan.Duration, &plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt)

		if scanErr != nil {
			common.Logger.Printf("Failed to scal plans rows at row (%d): %v\n", counter, scanErr)
		} else {
			plans = append(plans, plan)
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

func GetPlanByID(db *sql.DB, id int) (*Plan, error) {
	query := `SELECT title, description, price, duration, createdAt, updatedAt, deletedAt FROM Plan WHERE id = ?`

	plan := &Plan{}

	scanErr := db.QueryRow(query, id).Scan(&plan.Title, &plan.Description, &plan.Price, &plan.Duration, &plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt)

	if scanErr != nil {
		return nil, fmt.Errorf("Failed to get a plan by ID: %wn", scanErr)
	}

	return plan, nil
}

func CreatePlan(db *sql.DB, plan Plan) (int64, error) {
	query := `
  INSERT INTO Plan (title, description, price, duration) VALUES (?, ?, ?, ?);
  `
	res, queryErr := db.Exec(query, plan.Title, plan.Description, plan.Price, plan.Duration)
	id, lastInsertErr := res.LastInsertId()

	if queryErr != nil || lastInsertErr != nil {
		return 0, fmt.Errorf("Failed to create a plan: %w\n", queryErr)
	}

	return id, nil
}

func DeletePlanByID(db *sql.DB, id int) error {
	query := `DELETE FROM Plan WHERE id = ?`

	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete a plan by ID: %w\n", err)
	}

	return nil
}

func ReplacePlan(db *sql.DB, plan Plan) error {
	tx, txErr := db.Begin()

	if tx != nil {
		return fmt.Errorf("Could not begin a new transaction to replace a plan: %w\n", txErr)
	}

	deleteQuery := `DELETE FROM Plan WHERE id = ?`

	_, deleteErr := tx.Exec(deleteQuery, plan.ID)

	if deleteErr != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete a plan by ID in a transaction: %w\n", deleteErr)
	}

	insertQuery := `INSERT INTO Plan (title, description, price, duration) VALUES (?, ?, ?, ?);`

	_, insertErr := tx.Exec(insertQuery, plan.Title, plan.Description, plan.Price, plan.Duration)

	if insertErr != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create a plan in a transaction: %w\n", insertErr)
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		tx.Rollback()
		return fmt.Errorf("Could not commit plan replacement: %w\n", commitErr)
	}

	return nil
}

func GetProducts(db *sql.DB) ([]Product, error) {
	query := `SELECT id, name, description, price, marka, categoryId FROM Product`

	rows, queryErr := db.Query(query)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Product{}, nil
		}

		return nil, fmt.Errorf("Failed to get all products: %w\n", queryErr)
	}

	defer rows.Close()

	products := []Product{}
	counter := 0

	for rows.Next() {
		product := Product{}

		scanErr := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Marka, &product.CategoryID)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a product from rows at row (%d): %v\n", counter, scanErr)
		} else {
			products = append(products, product)
		}

		counter++
	}

	return products, nil
}

// / Categories come with empty arrays
func GetProductsWithCategories(db *sql.DB) ([]Product, error) {
	query := `SELECT P.id, P.name, P.description, P.price, P.marka, P.categoryId, C.id, C.name FROM Product AS P LEFT JOIN ProductCategory AS C`

	rows, queryErr := db.Query(query)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Product{}, nil
		}

		return nil, fmt.Errorf("Failed to get all products: %w\n", queryErr)
	}

	defer rows.Close()

	products := []Product{}
	counter := 0

	for rows.Next() {
		product := Product{Category: &ProductCategory{}}

		scanErr := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Marka, &product.Category.ID, &product.CategoryID, &product.Category.Name)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a product from rows at row (%d): %v\n", counter, scanErr)
		} else {
			products = append(products, product)
		}

		counter++
	}

	return products, nil
}

func GetProductByID(db *sql.DB, id int) (*Product, error) {
	query := `SELECT P.id, P.name, P.description, P.price, P.marka, P.categoryId FROM Product AS P WHERE id = ?`

	product := &Product{}

	scanErr := db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Marka, &product.CategoryID)

	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, fmt.Errorf("Product with ID (%d) not found\n", id)
		}

		return nil, fmt.Errorf("Failed to get a product by ID (%d): %w\n", id, scanErr)
	}

	return product, nil
}

// Category comes with an empty array
func GetProductWithCategoryByID(db *sql.DB, id int) (*Product, error) {
	query := `SELECT P.id, P.name, P.description, P.price, P.marka, P.categoryId, C.id, C.name FROM Product AS P LEFT JOIN ProductCategory WHERE P.id = ?`

	product := &Product{Category: &ProductCategory{}}

	scanErr := db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Marka, &product.CategoryID, &product.Category.ID, &product.Category.Name)

	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, fmt.Errorf("Product with ID (%d) not found\n", id)
		}

		return nil, fmt.Errorf("Failed to get a product by ID (%d): %w\n", id, scanErr)

	}

	return product, nil
}

func GetLandingPageGeneralInfo(db *sql.DB) (*LandingPageGeneralData, error) {
	query := `SELECT title, starterSentence, secondStarterSentence, plansParagraph FROM LandingPageData LIMIT 1`

	info := &LandingPageGeneralData{}

	scanErr := db.QueryRow(query).Scan(&info.Title, &info.StarterSentence, &info.SecondStarterSentence, &info.PlansParagraph)

	if scanErr != nil {
		return nil, fmt.Errorf("Failed to get landing page general info: %w\n", scanErr)
	}

	return info, nil
}

func UpdateLandingPageGeneralInfo(db *sql.DB, info LandingPageGeneralData) error {
	query := `UPDATE LandingPageData SET title = ?, starterSentence = ?, secondStarterSentence = ?, plansParagraph = ?`

	_, execErr := db.Exec(query, info.Title, info.StarterSentence, info.SecondStarterSentence, info.PlansParagraph)

	if execErr != nil {
		return fmt.Errorf("Failed to update landing page general info: %w\n", execErr)
	}

	return nil
}

func GetPlansParagraph(db *sql.DB) (string, error) {
	query := `SELECT plansParagraph FROM LandingPageData`

	var text string

	scanErr := db.QueryRow(query).Scan(&text)

	if scanErr != nil {
		return "", fmt.Errorf("Failed to get plans paragraph: %w\n", scanErr)
	}

	return text, nil
}

func UpdatePlansParagraph(db *sql.DB, text string) error {
	query := `UPDATE LandingPageData SET plansParagraph = ?`

	_, execErr := db.Exec(query, text)
	if execErr != nil {
		return fmt.Errorf("Failed to update plans paragraph: %w\n", execErr)
	}

	return nil
}

func GetAdsInfo(db *sql.DB) (*AdsInfo, error) {
	query := `SELECT adsOnImageBoldText, adsOnImageDescription FROM LandingPageData`

	info := &AdsInfo{}

	scanErr := db.QueryRow(query).Scan(&info.Title, &info.Description)
	if scanErr != nil {
		return nil, fmt.Errorf("Failed to get ads info: %w\n", scanErr)
	}

	return info, nil
}

func UpdateAdsInfo(db *sql.DB, info AdsInfo) error {
	query := `UPDATE LandingPageData SET adsOnImageBoldText = ?, adsOnImageDescription = ?`

	_, execErr := db.Exec(query, info.Title, info.Description)
	if execErr != nil {
		return fmt.Errorf("Failed to update ads info: %w\n", execErr)
	}

	return nil
}

func GetContacts(db *sql.DB) (*Contacts, error) {
	query := `SELECT emailContact, twitterContact, instigramContact, whatsappContact, facebookContact FROM LandingPageData`

	contacts := &Contacts{}

	scanErr := db.QueryRow(query).Scan(&contacts.Email, &contacts.Twitter, &contacts.Instagram, &contacts.WhatsApp, &contacts.Facebook)
	if scanErr != nil {
		return nil, fmt.Errorf("Failed to get contacts: %w\n", scanErr)
	}

	return contacts, nil
}

func UpdateContacts(db *sql.DB, contacts Contacts) error {
	query := `UPDATE LandingPageData SET emailContact = ?, twitterContact = ?, instigramContact = ?, whatsappContact = ?, facebookContact = ?`

	_, execErr := db.Exec(query, contacts.Email, contacts.Twitter, contacts.Instagram, contacts.WhatsApp, contacts.Facebook)
	if execErr != nil {
		return fmt.Errorf("Failed to update contacts: %w\n", execErr)
	}

	return nil
}

func GetLandingPageInfo(db *sql.DB) (*LandingPageData, error) {
	query := `SELECT title, starterSentence, secondStarterSentence, plansParagraph, adsOnImageBoldText, adsOnImageDescription, emailContact, twitterContact, facebookContact, instigramContact, whatsappContact FROM LandingPageData LIMIT 1`

	info := &LandingPageData{}

	scanErr := db.QueryRow(query).Scan(
		&info.Title,
		&info.StarterSentence,
		&info.SecondStarterSentence,
		&info.PlansParagraph,
		&info.AdsOnImageBoldText,
		&info.AdsOnImageDescription,
		&info.EmailContact,
		&info.TwitterContact,
		&info.FacebookContact,
		&info.InstigramContact,
		&info.WhatsappContact)

	if scanErr != nil {
		return nil, fmt.Errorf("Failed to get landing page general info: %w\n", scanErr)
	}

	return info, nil
}

func GetAllEvents(db *sql.DB) ([]Event, error) {
	query := `SELECT id, event, target, actorId, targetId, date FROM Event`

	rows, queryErr := db.Query(query)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Event{}, nil
		}

		common.Logger.Printf("Failed to get all events: %v\n", queryErr)
	}

	defer rows.Close()

	events := []Event{}
	counter := 0
	for rows.Next() {
		event := Event{Actor: nil}

		scanErr := rows.Scan(&event.ID, &event.Event, &event.Target, &event.ActorID, &event.TargetID, &event.Date)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan event from rows at row (%d): %v\n", counter, scanErr)
		} else {
			events = append(events, event)
		}

		counter++
	}

	return events, nil
}

func DidUserSeeEvent(db *sql.DB, userId, eventId int64) (bool, error) {
	query := `SELECT EXISTS(
    SELECT 1 FROM SeenEvent
    WHERE eventId = ? AND userId = ?
  );`

	var seen bool

	scanErr := db.QueryRow(query, eventId, userId).Scan(&seen)
	if scanErr != nil {
		return false, fmt.Errorf("Failed to check if an event was seen by a user: %w\n", scanErr)
	}

	return seen, nil
}

func MarkEventAsSeen(db *sql.DB, userId, eventId int64) error {
	query := `INSERT IGNORE INTO SeenEvent (eventId, userId) VALUES (?, ?)`

	_, execErr := db.Exec(query, eventId, userId)
	if execErr != nil {
		return fmt.Errorf("Failed to mark an event as seen: %w\n", execErr)
	}

	return nil
}

func MarkAllEventsAsSeen(db *sql.DB, userId int64) error {
	unreadQuery := `SELECT E.id FROM Event AS E WHERE NOT EXISTS (
    SELECT 1 FROM SeenEvent AS S WHERE S.eventId = E.Id AND S.userId = ?
  )`

	rows, queryErr := db.Query(unreadQuery, userId)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return nil
		}

		return fmt.Errorf("Failed to mark all events as seen by a user (querying unread): %w\n", queryErr)
	}

	defer rows.Close()

	counter := 0
	for rows.Next() {
		var eventId int

		scanErr := rows.Scan(&eventId)

		if scanErr != nil {
			common.Logger.Printf("Failed to scan an unseen event rows at row (%d): %v\b", counter, scanErr)
		} else {
			_, execErr := db.Exec(`INSERT IGNORE INTO SeenEvent (eventId, userrId) VALUES (?, ?)`, eventId, userId)
			if execErr != nil {
				common.Logger.Printf("Failed to mark an event (id: %d) as seen by a user: %v\n", eventId, execErr)
			}
		}

		counter++
	}

	return nil
}

func GetAllExercises(db *sql.DB) ([]Excercise, error) {
	query := `SELECT id, name, description, categoryId FROM Excercise`

	rows, queryErr := db.Query(query)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Excercise{}, nil
		}

		return nil, fmt.Errorf("Failed to get all exercises: %w\n", queryErr)
	}

	defer rows.Close()

	exercises := []Excercise{}
	counter := 0

	for rows.Next() {
		exercise := Excercise{}

		scanErr := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Description, &exercise.CategoryID)

		if scanErr != nil {
			common.Logger.Printf("Failed to scan an exercise from rows at row (%d): %v\n", counter, scanErr)
		} else {
			exercises = append(exercises, exercise)
		}

		counter++
	}

	return exercises, nil
}

func GetAllExercisesOfSection(db *sql.DB, sectionId int64) ([]Excercise, error) {
	query := `SELECT id, name, description FROM Excercise WHERE categoryId = ?`

	rows, queryErr := db.Query(query, sectionId)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Excercise{}, nil
		}

		return nil, fmt.Errorf("Failed to get all exercises: %w\n", queryErr)
	}

	defer rows.Close()

	exercises := []Excercise{}
	counter := 0

	for rows.Next() {
		exercise := Excercise{CategoryID: sectionId}

		scanErr := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Description)

		if scanErr != nil {
			common.Logger.Printf("Failed to scan an exercise from rows at row (%d): %v\n", counter, scanErr)
		} else {
			exercises = append(exercises, exercise)
		}

		counter++
	}

	return exercises, nil
}

func GetAllExercisesWithSections(db *sql.DB) ([]Excercise, error) {
	query := `SELECT E.id, E.name, E.description, E.categoryId, S.id, S.Name FROM Excercise AS E LEFT JOIN ExcerciseCategory AS C`

	rows, queryErr := db.Query(query)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Excercise{}, nil
		}

		return nil, fmt.Errorf("Failed to get all exercises: %w\n", queryErr)
	}

	defer rows.Close()

	exercises := []Excercise{}
	counter := 0

	for rows.Next() {
		exercise := Excercise{Category: &ExcerciseCategory{}}

		scanErr := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Description, &exercise.CategoryID, &exercise.Category.ID, &exercise.Category.Name)

		if scanErr != nil {
			common.Logger.Printf("Failed to scan an exercise from rows at row (%d): %v\n", counter, scanErr)
		} else {
			exercises = append(exercises, exercise)
		}

		counter++
	}

	return exercises, nil
}

func GetAllExerciseSections(db *sql.DB) ([]ExcerciseCategory, error) {
	query := `SELECT id, name FROM ExcerciseCategory`

	rows, queryErr := db.Query(query)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []ExcerciseCategory{}, nil
		}

		return nil, fmt.Errorf("Failed to get all exercise sections: %w\n", queryErr)
	}

	defer rows.Close()

	sections := []ExcerciseCategory{}
	counter := 0
	for rows.Next() {
		section := ExcerciseCategory{}

		scanErr := rows.Scan(&section.ID, &section.Name)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan exercise section from rows at row (%d): %v\n", counter, scanErr)
		} else {
			sections = append(sections, section)
		}

		counter++
	}

	return sections, nil
}

func GetAllExerciseSectionsWithExercises(db *sql.DB) ([]ExcerciseCategory, error) {
	query := `SELECT id, name FROM ExcerciseCategory`

	rows, queryErr := db.Query(query)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []ExcerciseCategory{}, nil
		}

		return nil, fmt.Errorf("Failed to get all exercise sections: %w\n", queryErr)
	}

	defer rows.Close()

	sections := []ExcerciseCategory{}
	counter := 0
	for rows.Next() {
		section := ExcerciseCategory{}

		scanErr := rows.Scan(&section.ID, &section.Name)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan exercise section from rows at row (%d): %v\n", counter, scanErr)
		} else {
			exercises, exerQueryErr := GetAllExercisesOfSection(db, section.ID)
			if exerQueryErr != nil {
				common.Logger.Printf("Failed to get all exercises of section (id: %d) of row (%d): %v\n", section.ID, counter, exerQueryErr)
			} else {
				section.Excercises = exercises
			}

			sections = append(sections, section)
		}

		counter++
	}

	return sections, nil
}

func GetExerciseSectionByIDWithExercises(db *sql.DB, id int64) (*ExcerciseCategory, error) {
	query := `SELECT name FROM ExcerciseCategory WHERE id = ?`

	row := db.QueryRow(query, id)

	section := ExcerciseCategory{ID: id}

	scanErr := row.Scan(&section.Name)

	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to get an exercise section by id with exercises (not found)")
		}

		return nil, fmt.Errorf("Failed to get an exercise section: %w\n", scanErr)
	}

	exercises, exerQueryErr := GetAllExercisesOfSection(db, id)
	if exerQueryErr == nil {
		section.Excercises = exercises
	}

	return &section, nil
}

func GetExerciseSectionByNameWithExercises(db *sql.DB, name string) (*ExcerciseCategory, error) {
	query := `SELECT id FROM ExcerciseCategory WHERE name = ? LIMIT 1`

	row := db.QueryRow(query, name)

	section := ExcerciseCategory{Name: name}

	scanErr := row.Scan(&section.ID)

	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to get an exercise section by name with exercises (not found)")
		}

		return nil, fmt.Errorf("Failed to get an exercise section: %w\n", scanErr)
	}

	exercises, exerQueryErr := GetAllExercisesOfSection(db, section.ID)
	if exerQueryErr == nil {
		section.Excercises = exercises
	}

	return &section, nil
}

func CreateExerciseSection(db *sql.DB, name string) (int64, error) {
	query := `INSERT INTO ExcerciseCategory (name) VALUES (?)`

	res, execErr := db.Exec(query, name)

	if execErr != nil {
		return 0, fmt.Errorf("Failed to create an exercise section: %w\n", execErr)
	}

	id, lastInsertErr := res.LastInsertId()
	if lastInsertErr != nil {
		return 0, fmt.Errorf("Failed to retrieve the inserted ID: %w\n", lastInsertErr)
	}

	return id, nil
}

func DeleteExerciseDeleteByName(db *sql.DB, name string) error {
	query := `DELETE FROM ExcerciseCategory WHERE name = ?`

	_, execErr := db.Exec(query, name)
	if execErr != nil {
		return fmt.Errorf("Failed to delete an exercise section (name: %s): %w\n", name, execErr)
	}

	return nil
}

func UpdateExerciseSectionByID(db *sql.DB, data ExcerciseCategory) error {
	query := `UPDATE ExcerciseCategory SET name = ? WHERE id = ?`

	_, execErr := db.Exec(query, data.Name, data.ID)

	if execErr != nil {
		return fmt.Errorf("Failed to update an exercise section by ID (id: %d): %w\n", data.ID, execErr)
	}

	return nil
}

func DeleteExerciseSectionByIDWithExercises(db *sql.DB, id int64) error {
	deleteSectionQuery := `DELETE FROM ExcerciseCategory WHERE id = ?`
	deleteExercisesQuery := `DELETE FROM Excercise WHERE categoryId = ?`

	tx, txErr := db.Begin()

	if txErr != nil {
		return fmt.Errorf("Failed to delete exercise section with exercises (failed to begine transaction): %w\n", txErr)
	}

	_, deleteExerErr := tx.Exec(deleteExercisesQuery, id)
	if deleteExerErr != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete exercise section with exercises (failed to delete exercises): %w\n", deleteExerErr)
	}

	_, deleteSectionErr := tx.Exec(deleteSectionQuery, id)
	if deleteSectionErr != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete exercise section with exercises (failed to delete section): %w\n", deleteSectionErr)
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete exercise section with exercises (failed to commit transaction): %w\n", commitErr)
	}

	return nil
}

func CountExercisesOfExerciseSectionByName(db *sql.DB, name string) (int, error) {
	query := `
    SELECT COUNT(*) 
    FROM Excercise AS E 
    WHERE EXISTS (
      SELECT id FROM ExcerciseCategory WHERE name = ? AND id = E.categoryId LIMIT 1
    )
  `
	var count int

	row := db.QueryRow(query, name)

	scanErr := row.Scan(&count)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return 0, fmt.Errorf("Failed to count exercises of an exercise section by name: (exercise section not found)")
		}
		return 0, fmt.Errorf("Failed to count exercises of an exercise section by name: %w\n", scanErr)
	}

	return count, nil
}

func CreateExercise(db *sql.DB, exercise Excercise) (int64, error) {
	query := `INSERT INTO Excercise (name, description, categoryId) VALUES(?, ?, ?)`

	res, execErr := db.Exec(query, exercise.Name, exercise.Description, exercise.CategoryID)
	if execErr != nil {
		return 0, fmt.Errorf("Failed to create exercise: %w\n", execErr)
	}

	id, scanErr := res.LastInsertId()
	if scanErr != nil {
		return 0, fmt.Errorf("Failed to obtain created exercise id: %w\n", scanErr)
	}

	return id, nil
}

func DeleteExerciseByID(db *sql.DB, id int64) error {
	query := `DELETE FROM Excercise WHERE id = ?`

	_, execErr := db.Exec(query, id)
	if execErr != nil {
		return fmt.Errorf("Failed to delete an exercise by ID (id: %d): %w\n", id, execErr)
	}

	return nil
}

func DeleteExerciseByName(db *sql.DB, name string) error {
	query := `DELETE FROM Excercise WHERE name = ?`

	_, execErr := db.Exec(query, name)
	if execErr != nil {
		return fmt.Errorf("Failed to delete an exercise by name (name: %s): %w\n", name, execErr)
	}

	return nil
}

func UpdateExercise(db *sql.DB, exercise Excercise) error {
	query := `UPDATE Excercise SET name = ?, description = ?, categoryId = ? WHERE id = ?`

	_, execErr := db.Exec(query, exercise.Name, exercise.Description, exercise.CategoryID, exercise.ID)
	if execErr != nil {
		return fmt.Errorf("Failed to update exercise: %w\n", execErr)
	}

	return nil
}

// Does not include users or subscribers
func GetAllComments(db *sql.DB, size, offset int) ([]SubscriberComment, error) {
	query := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL`
	limitedQuery := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL LIMIT ? OFFSET ?`

	var (
		rows     *sql.Rows
		queryErr error
	)

	if size == 0 {
		rows, queryErr = db.Query(query)
	} else {
		rows, queryErr = db.Query(limitedQuery, size, offset)
	}

	if queryErr != nil {
		return nil, fmt.Errorf("Failed to get all, undeleted comments: %w\n", queryErr)
	}

	defer rows.Close()

	comments := []SubscriberComment{}
	counter := 0

	for rows.Next() {
		comment := SubscriberComment{}

		scanErr := rows.Scan(&comment.ID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt, &comment.SenderID, &comment.SubscriberID)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a comment from rows at row (%d): %v\n", counter, scanErr)
		} else {
			comments = append(comments, comment)
		}

		counter++
	}

	return comments, nil
}

func GetAllCommentsIncludes(db *sql.DB, size, offset int, includeUsers, includeSubscribers bool) ([]SubscriberComment, error) {
	query := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL`
	limitedQuery := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL LIMIT ? OFFSET ?`

	var (
		rows     *sql.Rows
		queryErr error
	)

	if size == 0 {
		rows, queryErr = db.Query(query)
	} else {
		rows, queryErr = db.Query(limitedQuery, size, offset)
	}

	if queryErr != nil {
		return nil, fmt.Errorf("Failed to get all, undeleted comments: %w\n", queryErr)
	}

	defer rows.Close()

	comments := []SubscriberComment{}
	counter := 0

	for rows.Next() {
		comment := SubscriberComment{}

		scanErr := rows.Scan(&comment.ID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt, &comment.SenderID, &comment.SubscriberID)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a comment from rows at row (%d): %v\n", counter, scanErr)
		} else {
			comments = append(comments, comment)

			if includeUsers {
				user, userQueryErr := GetUserByID(db, int64(comment.SenderID))
				if userQueryErr != nil {
					common.Logger.Printf("Failed to get a user for a comment (relation) at row (%d): %v\n", counter, userQueryErr)
				} else {
					comment.Sender = user
				}
			}

			if includeSubscribers {
				sub, subQueryErr := GetSubscriberByID(db, comment.SubscriberID)
				if subQueryErr != nil {
					common.Logger.Printf("Failed to get a subscriber for a comment (relation) at row (%d): %v\n", counter, subQueryErr)
				} else {
					comment.Subscriber = sub
				}
			}
		}

		counter++
	}

	return comments, nil
}

func GetAllCommentsOfUserID(db *sql.DB, id int64, size, offset int) ([]SubscriberComment, error) {
	query := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL AND senderId = ?`
	limitedQuery := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL AND senderId = ? LIMIT ? OFFSET ?`

	var (
		rows     *sql.Rows
		queryErr error
	)

	if size == 0 {
		rows, queryErr = db.Query(query, id)
	} else {
		rows, queryErr = db.Query(limitedQuery, id, size, offset)
	}

	if queryErr != nil {
		return nil, fmt.Errorf("Failed to get all, undeleted comments: %w\n", queryErr)
	}

	defer rows.Close()

	comments := []SubscriberComment{}
	counter := 0

	for rows.Next() {
		comment := SubscriberComment{}

		scanErr := rows.Scan(&comment.ID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt, &comment.SenderID, &comment.SubscriberID)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a comment from rows at row (%d): %v\n", counter, scanErr)
		} else {
			comments = append(comments, comment)
		}

		counter++
	}

	return comments, nil
}

func GetAllCommentsOfSubscriberID(db *sql.DB, id int64, size, offset int) ([]SubscriberComment, error) {
	query := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL AND subscriberId = ?`
	limitedQuery := `SELECT id, text, createdAt, updatedAt, senderId, subscriberId FROM SubscriberComment WHERE deletedAt IS NULL AND subscriberId = ? LIMIT ? OFFSET ?`

	var (
		rows     *sql.Rows
		queryErr error
	)

	if size == 0 {
		rows, queryErr = db.Query(query, id)
	} else {
		rows, queryErr = db.Query(limitedQuery, id, size, offset)
	}

	if queryErr != nil {
		return nil, fmt.Errorf("Failed to get all, undeleted comments: %w\n", queryErr)
	}

	defer rows.Close()

	comments := []SubscriberComment{}
	counter := 0

	for rows.Next() {
		comment := SubscriberComment{}

		scanErr := rows.Scan(&comment.ID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt, &comment.SenderID, &comment.SubscriberID)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a comment from rows at row (%d): %v\n", counter, scanErr)
		} else {
			comments = append(comments, comment)
		}

		counter++
	}

	return comments, nil
}

func CreateComment(db *sql.DB, comment SubscriberComment) (int64, error) {
	query := `INSERT INTO SubscriberComment (text, senderId, subscriberId) VALUES (?, ?)`
	res, execErr := db.Exec(query, comment.Text, comment.SenderID, comment.SubscriberID)

	if execErr != nil {
		return 0, fmt.Errorf("Failed to create a comment: %w\n", execErr)
	}

	id, lastInsertErr := res.LastInsertId()
	if lastInsertErr != nil {
		return 0, fmt.Errorf("Failed to retrieve last inserted ID: %w\n", lastInsertErr)
	}

	return id, nil
}

func DeleteCommentByID(db *sql.DB, id int64) error {
	query := `DELETE FROM SubscriberComment WHERE id = ?`

	_, execErr := db.Exec(query, id)
	if execErr != nil {
		return fmt.Errorf("Failed to delete a comment by ID (id: %d): %w\n", id, execErr)
	}

	return nil
}

func GetAllAnnouncements(db *sql.DB) ([]Message, error) {
	query := `SELECT M.id, M.text, M.sent FROM Message AS M`

	rows, queryErr := db.Query(query)
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Message{}, nil
		}
		return nil, fmt.Errorf("Failed to get all messages: %w\n", queryErr)
	}

	defer rows.Close()

	messages := make([]Message, 0, 1)
	counter := 0

	for rows.Next() {
		message := Message{}

		scanErr := rows.Scan(&message.ID, &message.Text, &message.Sent)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a message from rows at a row (%d): %v\n", counter, scanErr)
		} else {
			messages = append(messages, message)
		}

		counter++
	}

	return messages, nil
}

func CreateAnnouncementToAll(db *sql.DB, text string) (int64, error) {
	userQuery := `SELECT id FROM User WHERE deletedAt IS NULL`
	annQuery := `INSERT INTO Message (text, sent) VALUES (?, CURRENT_TIMESTAMP)`
	readQuery := `INSERT INTO MessageRead (userId, messageId, read) VALUES (?, ?, 0)`

	tx, txErr := db.Begin()
	if txErr != nil {
		return 0, fmt.Errorf("Failed to create an announcement to all (failed transaction): %w\n", txErr)
	}

	res, annErr := tx.Exec(annQuery, text)
	if annErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at announcement creation: %w\n", rollErr)
		}

		return 0, fmt.Errorf("Failed to create announcement: %w\n", annErr)
	}

	insertedId, idErr := res.LastInsertId()
	if idErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at announcement creation id retrieval: %w\n", rollErr)
		}

		return 0, fmt.Errorf("Failed to retrieve created announcement ID: %w\n", idErr)
	}

	userRows, userQueryErr := tx.Query(userQuery)
	if userQueryErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at user ids fetching: %w\n", rollErr)
		}

		return 0, fmt.Errorf("Failed to get user IDs: %w\n", userQueryErr)
	}

	defer userRows.Close()

	counter := 0

	for userRows.Next() {
		var id int64

		idScanErr := userRows.Scan(&id)
		if idScanErr != nil {
			common.Logger.Printf("Failed to scan user ID from rows at row (%d): %v\n", counter, idScanErr)
		} else {
			_, execErr := tx.Exec(readQuery, id, insertedId)
			if execErr != nil {
				rollErr := tx.Rollback()
				if rollErr != nil {
					return 0, fmt.Errorf("Failed to rollback on transaction at message read creation: %w\n", rollErr)
				}

				return 0, fmt.Errorf("Failed to create message read status for all users at user ID (id: %d): %w\n", id, execErr)
			}
		}

		counter++
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at committing transaction: %w\n", rollErr)
		}
		return 0, fmt.Errorf("Failed to create announcements (failed to commit transaction): %w\n", commitErr)
	}

	return insertedId, nil
}

func CreateAnnouncementToUserIDs(db *sql.DB, text string, ids ...int64) (int64, error) {
	annQuery := `INSERT INTO Message (text, sent) VALUES (?, CURRENT_TIMESTAMP)`
	readQuery := `INSERT INTO MessageRead (userId, messageId, read) VALUES (?, ?, 0)`

	tx, txErr := db.Begin()
	if txErr != nil {
		return 0, fmt.Errorf("Failed to create an announcement to all (failed transaction): %w\n", txErr)
	}

	res, annErr := tx.Exec(annQuery, text)
	if annErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at announcement creation: %w\n", rollErr)
		}

		return 0, fmt.Errorf("Failed to create announcement: %w\n", annErr)
	}

	insertedId, idErr := res.LastInsertId()
	if idErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at announcement creation id retrieval: %w\n", rollErr)
		}

		return 0, fmt.Errorf("Failed to retrieve created announcement ID: %w\n", idErr)
	}

	counter := 0

	for _, id := range ids {
		_, execErr := tx.Exec(readQuery, id, insertedId)
		if execErr != nil {
			rollErr := tx.Rollback()
			if rollErr != nil {
				return 0, fmt.Errorf("Failed to rollback on transaction at message read creation: %w\n", rollErr)
			}

			return 0, fmt.Errorf("Failed to create message read status for all users at user ID (id: %d): %w\n", id, execErr)
		}

		counter++
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at committing transaction: %w\n", rollErr)
		}
		return 0, fmt.Errorf("Failed to create announcements (failed to commit transaction): %w\n", commitErr)
	}

	return insertedId, nil
}

func MarkMessageAsRead(db *sql.DB, userId, messageId int64) error {
	query := `UPDATE MessageRead SET read = 1 WHERE userId = ? AND messageId = ?`

	_, execErr := db.Exec(query, userId, messageId)
	if execErr != nil {
		return fmt.Errorf("Failed to mark a message as read: %w\n", execErr)
	}

	return nil
}

func GetAllTrainers(db *sql.DB) ([]Trainer, error) {
	query := `SELECT id, name, job, description, intsigram, facebook, twitter FROM Trainer`

	rows, queryErr := db.Query(query)
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return []Trainer{}, nil
		}
		return nil, fmt.Errorf("Failed to get all trainers: %W\n", queryErr)
	}

	defer rows.Close()

	trainers := []Trainer{}
	counter := 0

	for rows.Next() {
		trainer := Trainer{}

		scanErr := rows.Scan(&trainer.ID, &trainer.Name, &trainer.Job, &trainer.Description, &trainer.Instigram, &trainer.Facebook, &trainer.Twitter)
		if scanErr != nil {
			common.Logger.Printf("Failed to scan a trainer from rows at row (%d): %v\n", counter, scanErr)
		} else {
			trainers = append(trainers, trainer)
		}

		counter++
	}

	return trainers, nil
}

func CreateTrainer(db *sql.DB, data Trainer) (int64, error) {
	query := `INSERT INTO Trainer (name, job, description, instigram, facebook, twitter) VALUES (?, ?, ?, ?, ?)`

	tx, txErr := db.Begin()
	if txErr != nil {
		return 0, fmt.Errorf("Failed to create a trainer (failed to begin transaction): %w\n", txErr)
	}

	res, execErr := tx.Exec(query, data.Name, data.Job, data.Description, data.Instigram, data.Facebook, data.Twitter)
	if execErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to create a trainer: (Failed to rollback on transaction): %w\n", rollErr)
		}

		return 0, fmt.Errorf("Failed to create a trainer: %w\n", execErr)
	}

	id, idErr := res.LastInsertId()
	if idErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to create a trainer: (Failed to rollback on transaction): %w\n", rollErr)
		}

		return 0, fmt.Errorf("Failed to retrieve created trainer ID: %w\n", idErr)
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, fmt.Errorf("Failed to rollback on transaction at committing transaction: %w\n", rollErr)
		}
		return 0, fmt.Errorf("Failed to create a trainer (failed to commit transaction): %w\n", commitErr)
	}
	return id, nil
}

func UpdateTrainer(db *sql.DB, data Trainer) error {
	query := `UPDATE Trainer SET name = ?, job = ?, description = ?, instigram = ?, facebook = ?, twitter = ? WHERE id = ?`

	_, execErr := db.Exec(query, data.Name, data.Job, data.Description, data.Instigram, data.Facebook, data.Twitter, data.ID)
	if execErr != nil {
		return fmt.Errorf("Failed to update a trainer (id: %d): %w\n", data.ID, execErr)
	}

	return nil
}

func DeleteTrainerByID(db *sql.DB, id int64) error {
	query := `DELETE FROM Trainer WHERE id = ?`
	_, execErr := db.Exec(query, id)
	if execErr != nil {
		return fmt.Errorf("Failed to delete a trainer by ID (id: %d): %w\n", id, execErr)
	}

	return nil
}
