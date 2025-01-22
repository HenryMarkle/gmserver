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
		common.Logger.Printf("request from origin %s\n", ctx.Request.Header.Get("Origin"))

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Accept-Encoding, Connection, Content-Length, Authorization, Cookie")
		// ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
			return
		}

		ctx.Next()
	}
}

func DynamicCORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		common.Logger.Printf("request from origin %s\n", origin)

		if origin != "" {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")

		} else {
			ctx.Header("Access-Control-Allow-Origin", "*")
		}
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Accept-Encoding, Connection, Content-Length, Authorization, Cookie")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}

		ctx.Next()
	}
}
