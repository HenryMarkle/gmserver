package api

import (
	"net/http"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AddUser(ctx *gin.Context) {
	data := dto.CreateUser_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v\n", bindErr)
		return
	}

	if exists := db.UserExistsByEmail(db.DB, data.Email); exists {
		ctx.String(http.StatusBadRequest, "Email address is already used")
		return
	}

	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if hashErr != nil {
		common.Logger.Printf("Failed to hash password: %v\n", hashErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	queryErr := db.AddAccount(db.DB, db.User{
		Email:     data.Email,
		Name:      data.Name,
		StartDate: data.StartDate,
		Password:  string(hashed),
		Session:   uuid.New().String(),
		Gender:    data.Gender,
		Salary:    data.Salary,
	})

	if queryErr != nil {
		common.Logger.Printf("Failed to create a new account: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}
func IsUserSignedIn(ctx *gin.Context) {}
func GetCurrentUserId(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userPtr := user.(*db.User)
	ctx.JSON(http.StatusOK, userPtr.ID)
}

func GetTotalSalaries(ctx *gin.Context) {
	sum, queryErr := db.GetTotalSalaries(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get total salaries: %v\n", queryErr)
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, sum)
}

func UpdateUser(ctx *gin.Context) {
	data := dto.UpdateUser_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v\n", bindErr)
		return
	}

	queryErr := db.UpdateUser(db.DB, db.User{
		ID:      data.ID,
		Email:   data.Email,
		Name:    data.Name,
		GymName: data.GymName,
		Gender:  data.Gender,
		Salary:  data.Salary,
		Age:     data.Age,
	})

	if queryErr != nil {
		common.Logger.Printf("Failed to update user: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func ChangeUserName(ctx *gin.Context)
func DeleteUser(ctx *gin.Context)
func DeleteUserById(ctx *gin.Context) {
	idStr := ctx.Query("id")
	permanentStr := ctx.Query("permanent")

	permanent, permaConvErr := strconv.ParseBool(permanentStr)
	id, idConvErr := strconv.ParseInt(idStr, 10, 64)

	if permaConvErr != nil && idConvErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid URL query parameters: id, permanent")
		return
	}

	if permaConvErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid URL query parameters: permanent")
		return
	}

	if idConvErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid URL query parameters: id")
		return
	}

	var queryErr error

	if permanent {
		queryErr = db.DeleteUserByID(db.DB, id)
	} else {
		queryErr = db.MarkUserAsDeleted(db.DB, id)
	}

	if queryErr != nil {
		common.Logger.Printf("Failed to delete a user: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func CountUsers(ctx *gin.Context) {
	count, queryErr := db.CountUsers(db.DB)

	if queryErr != nil {
		common.Logger.Printf("Failed to count users: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

func GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")

	if email == "" {
		ctx.String(http.StatusBadRequest, "Invalid URL parameter: email")
		return
	}

	user, queryErr := db.GetUserByEmail(db.DB, email)
	if queryErr != nil {
		common.Logger.Printf("Failed to get a user by email: %v\n", queryErr)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, dto.User_Res{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		GymName:    user.GymName,
		Gender:     user.Gender,
		LastLogin:  user.LastLogin,
		Permission: user.Permission,
		Salary:     user.Salary,
		Age:        user.Age,
	})
}
func GetUsersLeftChartData(ctx *gin.Context)
func GetUsersCreatedChartData(ctx *gin.Context)
func GetUserById(ctx *gin.Context) {
	idStr := ctx.Query("id")

	if idStr == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: id")
		return
	}

	id, idConvErr := strconv.ParseInt(idStr, 10, 64)
	if idConvErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid query parameter: id")
		return
	}

	user, queryErr := db.GetUserByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to get user by id: %v\n", queryErr)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, dto.User_Res{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		GymName:    user.GymName,
		Gender:     user.Gender,
		LastLogin:  user.LastLogin,
		Permission: user.Permission,
		Salary:     user.Salary,
		Age:        user.Age,
	})
}

func GetGymName(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userPtr := user.(*db.User)

	ctx.String(http.StatusOK, userPtr.GymName)
}

func GetCurrentUser(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userPtr := user.(*db.User)

	ctx.JSON(http.StatusOK, dto.User_Res{
		ID:         userPtr.ID,
		Email:      userPtr.Email,
		Name:       userPtr.Name,
		GymName:    userPtr.GymName,
		Gender:     userPtr.Gender,
		LastLogin:  userPtr.LastLogin,
		Permission: userPtr.Permission,
		Salary:     userPtr.Salary,
		Age:        userPtr.Age,
	})
}

func ChangeGymName(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userPtr := user.(*db.User)

	var newGymName string

	bindErr := ctx.ShouldBindJSON(&newGymName)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v\n", bindErr)
		return
	}

	queryErr := db.ChangeGymName(db.DB, userPtr.ID, newGymName)
	if queryErr != nil {
		common.Logger.Printf("Failed to update the gym name of a user: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetAllUsers(ctx *gin.Context) {
	users, queryErr := db.GetAllUsers(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get all users: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dtoUsers := make([]dto.User_Res, 0, len(users))
	for _, user := range users {
		dtoUsers = append(dtoUsers, dto.User_Res{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			GymName:    user.GymName,
			Gender:     user.Gender,
			LastLogin:  user.LastLogin,
			Permission: user.Permission,
			Salary:     user.Salary,
			Age:        user.Age,
		})
	}

	ctx.JSON(http.StatusOK, dtoUsers)
}

// All announcement for a user
func GetAllAnnouncments(ctx *gin.Context) {
	announcements, queryErr := db.GetAllAnnouncements(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get accouncements: %v\n", queryErr)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, announcements)
}

func CreateAnnouncement(ctx *gin.Context) {
	data := dto.CreateAnnouncement_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	var (
		id       int64
		queryErr error
	)

	if data.All {
		id, queryErr = db.CreateAnnouncementToAll(db.DB, data.Text)
	} else {
		id, queryErr = db.CreateAnnouncementToUserIDs(db.DB, data.Text, data.ToUsers...)
	}

	if queryErr != nil {
		common.Logger.Printf("Failed to create announcement (to all: %t): %v\n", data.All, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func MarkAsRead(ctx *gin.Context) {
	data := dto.MarkMessageAsRead_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.MarkMessageAsRead(db.DB, data.UserID, data.MessageID)
	if queryErr != nil {
		common.Logger.Printf("Failed to mark a message as read: %v\n", queryErr)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
