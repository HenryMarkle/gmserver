package api

import (
	"net/http"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/gin-gonic/gin"
)

func GetAllEvents(ctx *gin.Context) {
	events, queryErr := db.GetAllEvents(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get events: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func DidUserSeeEvent(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "userId", "eventId")

	if params == nil {
		return
	}

	userId, eventId := params[0], params[1]

	result, queryErr := db.DidUserSeeEvent(db.DB, userId, eventId)
	if queryErr != nil {
		common.Logger.Printf("Failed to get if a user had seen an event: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func MarkEventAsSeen(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userPtr, _ := user.(*db.User)

	params := NonEmptyQueryInt64OrAbort(ctx, "eventId")
	if params == nil {
		return
	}
	eventId := params[0]

	queryErr := db.MarkEventAsSeen(db.DB, userPtr.ID, eventId)
	if queryErr != nil {
		common.Logger.Printf("Failed to mark an event as seen by a user: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func MarkAllEventsAsSeen(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userPtr, _ := user.(*db.User)

	queryErr := db.MarkAllEventsAsSeen(db.DB, userPtr.ID)
	if queryErr != nil {
		common.Logger.Printf("Failed to mark all events as seen by a user: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
