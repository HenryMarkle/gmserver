package api

import (
	"net/http"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
	"github.com/gin-gonic/gin"
)

func CreateAdvice(ctx *gin.Context) {
	req := dto.CreateAdvice_Req{}

	bindErr := ctx.ShouldBindBodyWithJSON(&req)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreateAdvice(db.DB, req.Title, req.Description)
	if queryErr != nil {
		common.Logger.Printf("failed to create advice: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

func UpdateAdviceByID(ctx *gin.Context) {
	req := dto.UpdateAdvice_Req{}

	bindErr := ctx.ShouldBindBodyWithJSON(&req)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateAdviceByID(db.DB, db.Advice{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	})

	if queryErr != nil {
		common.Logger.Printf("failed to update advice: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetAdviceByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	advice, queryErr := db.GetAdviceByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("failed to get advice by id: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, advice)
}

func GetAllAdvice(ctx *gin.Context) {
	advices, queryErr := db.GetAllAdvice(db.DB)
	if queryErr != nil {
		common.Logger.Printf("failed to get all advices: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, advices)
}

func DeleteAdviceByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	queryErr := db.DeleteAdviceByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("failed to get delete advice by id: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
