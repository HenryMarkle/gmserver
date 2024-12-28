package api

import (
	"net/http"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/gin-gonic/gin"
)

func AuthCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authCookie, err := ctx.Cookie("gmserver-session")
		if err == nil {
			ctx.Set("session", authCookie)
		}

		ctx.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, cookieErr := ctx.Cookie("gmserver-session")
		if cookieErr != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, queryErr := db.GetUserBySession(db.DB, sessionId)
		if queryErr != nil {
			common.Logger.Printf("Middleware error [Auth]: %v\n", queryErr)
			ctx.AbortWithStatus(500)
			return
		}

		if user == nil {
			common.Logger.Printf("Middleware [Auth]: attempted to access secured API without authorization\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("user")
		if !exists {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userPtr, ok := user.(*db.User)
		if !ok {
			common.Logger.Println("Middleware error [AdminOnly]: invalid user type")
			ctx.AbortWithStatus(500)
			return
		}

		if userPtr.Permission != 1 {
			common.Logger.Println("Middleware [AdminOnly]: unauthorized access to an admin-only API")
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}

func AllowAllOrigins() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}
