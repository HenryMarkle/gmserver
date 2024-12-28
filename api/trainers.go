package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
)

func GetTrainers(ctx *gin.Context) {
	trainers, queryErr := db.GetAllTrainers(db.DB)
	if queryErr != nil {
		common.Logger.Printf("failed to get trainers: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, trainers)
}

func CreateTrainer(ctx *gin.Context) {
	data := dto.CreateTrainer_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreateTrainer(db.DB, db.Trainer{
		Name:        data.Name,
		Job:         data.Job,
		Description: data.Description,
		Instigram:   data.Facebook,
		Facebook:    data.Facebook,
		Twitter:     data.Twitter,
	})

	if queryErr != nil {
		common.Logger.Printf("Failed to create a trainer: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func ReplaceTrainerById(ctx *gin.Context) {
	data := dto.UpdateTrainer_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateTrainer(db.DB, db.Trainer{
		Name:        data.Name,
		Job:         data.Job,
		Description: data.Description,
		Instigram:   data.Facebook,
		Facebook:    data.Facebook,
		Twitter:     data.Twitter,
		ID:          data.ID,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to update a trainer: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func DeleteTrainerById(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "id")

	if params == nil {
		return
	}

	id := params[0]

	queryErr := db.DeleteTrainerByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a trainer: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
