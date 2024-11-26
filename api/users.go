package api

import (
	"net/http"

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
func DeleteUserById(ctx *gin.Context)
func CountUsers(ctx *gin.Context)
func GetUserByEmail(ctx *gin.Context)
func GetUsersLeftChartData(ctx *gin.Context)
func GetUsersCreatedChartData(ctx *gin.Context)
func GetUserById(ctx *gin.Context)
func GetGymName(ctx *gin.Context)
func GetCurrentUser(ctx *gin.Context)
func ChangeGymName(ctx *gin.Context)
func GetAllUsers(ctx *gin.Context)
func GetAllAnnouncments(ctx *gin.Context)
func CreateAnnouncement(ctx *gin.Context)
func MarkAsRead(ctx *gin.Context)
