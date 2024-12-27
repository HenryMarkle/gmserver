package api

import (
	"net/http"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
	"github.com/gin-gonic/gin"
)

func GetAllComments(ctx *gin.Context) {
	comments, queryErr := db.GetAllComments(db.DB, 0, 0)
	if queryErr != nil {
		common.Logger.Printf("Failed to get comments: %v\n", queryErr)
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func GetAllCommentsOfManager(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	if idStr == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: id")
	}

	id, convErr := strconv.ParseInt(idStr, 10, 64)

	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid query parameter: id: %v\n", convErr)
		return
	}

	comments, queryErr := db.GetAllCommentsOfUserID(db.DB, id, 0, 0)
	if queryErr != nil {
		common.Logger.Printf("Failed to get comments: %v\n", queryErr)
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func GetAllCommentsOfSubscriber(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	if idStr == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: id")
	}

	id, convErr := strconv.ParseInt(idStr, 10, 64)

	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid query parameter: id: %v\n", convErr)
		return
	}

	comments, queryErr := db.GetAllCommentsOfSubscriberID(db.DB, id, 0, 0)
	if queryErr != nil {
		common.Logger.Printf("Failed to get comments: %v\n", queryErr)
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func CreateComment(ctx *gin.Context) {
	data := dto.CreateComment_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreateComment(db.DB, db.SubscriberComment{
		Text:         data.Text,
		SenderID:     data.SenderID,
		SubscriberID: data.SubscriberID,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to create a comment: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func DeleteComment(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	if idStr == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: id")
	}

	id, convErr := strconv.ParseInt(idStr, 10, 64)

	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid query parameter: id")
		return
	}

	queryErr := db.DeleteCommentByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete comment: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
