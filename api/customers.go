package api

import (
	"net/http"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
	"github.com/gin-gonic/gin"
)

func GetTotalIncome(ctx *gin.Context) {
	total, queryErr := db.GetTotalSubscriberPaymentAmount(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get total income: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, total)
}

func CountCustomers(ctx *gin.Context) {
	count, queryErr := db.GetSubscriberCount(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to subscriber count: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

func CountCustomersEndingIn(ctx *gin.Context) {
	var time = ctx.Query("date")

	count, queryErr := db.GetAllSubscribersEndingBefore(db.DB, time)
	if queryErr != nil {
		common.Logger.Printf("Failed to count ending dubscribers: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

func CountCustomersExpiring(ctx *gin.Context) {
	count, queryErr := db.GetAllExpiredSubscribers(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to count ended dubscribers: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

func CreateCustomer(ctx *gin.Context) {
	data := dto.CreateSubscriber_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.CreateSubscriber(db.DB, db.Subscriber{
		Name:          data.Name,
		Surname:       data.Surname,
		StartedAt:     data.StartedAt,
		EndsAt:        data.EndsAt,
		Gender:        data.Gender,
		Age:           data.Age,
		PaymentAmount: data.PaymentAmount,
		BucketPrice:   data.BucketPrice,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to create customer: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetAllCustomers(ctx *gin.Context) {
	limitStr := ctx.Query("limit")

	var limit int64
	var convErr error

	if limitStr != "" {
		limit, convErr = strconv.ParseInt(limitStr, 10, 32)
		if convErr != nil {
			ctx.String(http.StatusBadRequest, "Invalid query parameter: limit")
			return
		}
	}

	subs, queryErr := db.GetAllSubscribers(db.DB, int(limit))
	if queryErr != nil {
		common.Logger.Printf("Failed to get subscribers: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, subs)
}

func GetCustomerByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid paramter: id")
		return
	}

	sub, queryErr := db.GetSubscriberByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to get a subscriber by ID: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, sub)
}

func DeleteCustomerByID(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "id")
	if params == nil {
		return
	}

	id := params[0]

	queryErr := db.DeleteSubscriberByID(db.DB, id, true)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a subscriber by ID: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func MarkCustomerAsDeleted(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "id")
	if params == nil {
		return
	}

	id := params[0]

	queryErr := db.DeleteSubscriberByID(db.DB, id, false)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a subscriber by ID: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateCustomerByID(ctx *gin.Context) {
	sub := db.Subscriber{}

	bindErr := ctx.ShouldBindJSON(&sub)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v\n", bindErr)
		return
	}

	queryErr := db.UpdateSubscriber(db.DB, sub)
	if queryErr != nil {
		common.Logger.Printf("Failed to update a subscriber by ID: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
