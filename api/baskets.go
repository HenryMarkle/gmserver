package api

import (
	"net/http"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/gin-gonic/gin"
)

func GetUserBasket(ctx *gin.Context) {
	userPtr, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userPtr.(*db.User)

	basket, queryErr := db.GetAllBasketProductsOfUser_WithProducts(db.DB, user.ID)
	if queryErr != nil {
		common.Logger.Printf("failed to get basket of user: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, basket)
}

func GetUserBasketByID(ctx *gin.Context) {
	userPtr, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userPtr.(*db.User)

	basketIdStr := ctx.Params.ByName("basketId")
	basketId, convErr := strconv.ParseInt(basketIdStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: basketId")
		return
	}

	basket, queryErr := db.GetProductBasketByID_WithProduct(db.DB, basketId)
	if queryErr != nil {
		common.Logger.Printf("failed to get basket of user: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if basket == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	if basket.CustomerID != user.ID {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, basket)
}

func AddToUserBasket(ctx *gin.Context) {
	userPtr, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userPtr.(*db.User)

	params := NonEmptyQueryInt64OrAbort(ctx, "productId", "quantity")

	if params == nil {
		return
	}

	productId, quantity := params[0], params[1]

	if quantity < 1 {
		ctx.String(http.StatusBadRequest, "Query parameter 'quantity' must be positive and non-zero.")
		return
	}

	createdId, queryErr := db.CreateProductBasket(db.DB, user.ID, productId, int(quantity))
	if queryErr != nil {
		common.Logger.Printf("failed to add product to basket: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, int(createdId))
}

func IncrementBasketQuantity(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "basketId")
	if params == nil {
		return
	}
	basketId := params[0]

	queryErr := db.IncrementBasketProductQuantityByID(db.DB, basketId)
	if queryErr != nil {
		common.Logger.Printf("failed to increment basket quantity: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func DecrementBasketQuantity(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "basketId")
	if params == nil {
		return
	}
	basketId := params[0]

	queryErr := db.DecrementBasketProductQuantityByID(db.DB, basketId)
	if queryErr != nil {
		common.Logger.Printf("failed to decrement basket quantity: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func DeleteBasket(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "basketId")
	if params == nil {
		return
	}
	basketId := params[0]

	queryErr := db.DeleteProductBasketByID(db.DB, basketId)
	if queryErr != nil {
		common.Logger.Printf("failed to delete basket: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
