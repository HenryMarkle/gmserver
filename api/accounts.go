package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
)

/*
* API handlers here
 */

func SignIn(ctx *gin.Context) {
	signin := dto.Signin_Req{}
	if bindErr := ctx.ShouldBindJSON(&signin); bindErr != nil {
		ctx.String(400, "Invalid request data: %v", bindErr)
		return
	}

	row := db.DB.QueryRow(`SELECT id, password FROM User WHERE email = ? AND deletedAt IS NULL`, signin.Email)

	var userId int
	var hashedPassword string

	scanErr := row.Scan(&userId, &hashedPassword)

	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			ctx.String(400, "Account not found")
			return
		}

		common.Logger.Printf("Failed to signin: %v\n", scanErr)
		ctx.Status(500)
		return
	}

	compErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(signin.Password))
	if compErr != nil {
		ctx.String(400, "Invalid crredentials")
		return
	}

	sessionId := uuid.New()

	_, execErr := db.DB.Exec(`UPDATE User SET lastLogin = CURRENT_TIMESTAMP, session = ? WHERE id = ?`, sessionId.String(), userId)
	if execErr != nil {
		common.Logger.Printf("Failed to sign in (db error): %v\n", execErr)
		ctx.Status(500)
		return
	}

	ctx.SetCookie("gmserver-session", sessionId.String(), 1000*60*60*6, "/", "", true, true)

	ctx.Status(200)
}

func Signout(ctx *gin.Context) {
	cookie, cookieErr := ctx.Cookie("gmserver-session")

	if cookieErr != nil {
		common.Logger.Printf("Could not signout: %v\n", cookieErr)
		ctx.Status(400)
		return
	}

	row := db.DB.QueryRow(`SELECT id FROM User WHERE session = ?`, cookie)

	var userId int
	scanErr := row.Scan(&userId)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			ctx.Status(200)
			return
		}

		common.Logger.Printf("Failed to signout: %v\n", scanErr)

		ctx.Status(500)
		return
	}

	_, execErr := db.DB.Exec(`UPDATE User SET session = '' WHERE id = ?`, userId)
	if execErr != nil {
		common.Logger.Printf("Failed to signout: %v\n", execErr)
		ctx.Status(500)
		return
	}

	ctx.SetCookie("gmserver-session", "", 0, "/", "", true, true)

	ctx.Status(200)
}

func ChangePassword(ctx *gin.Context) {
	data := dto.ChangePassword_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userPtr := user.(*db.User)

	compErr := bcrypt.CompareHashAndPassword([]byte(userPtr.Password), []byte(data.OldPassword))
	if compErr != nil {
		ctx.String(http.StatusBadRequest, "Incorrect credentials")
		return
	}

	queryErr := db.ChangeUserPassword(db.DB, userPtr.ID, data.NewPassword)
	if queryErr != nil {
		common.Logger.Printf("Failed to change user password: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
