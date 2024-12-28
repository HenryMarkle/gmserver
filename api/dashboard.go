package api

import (
	"net/http"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
	"github.com/gin-gonic/gin"
)

func GetHomeInfo(ctx *gin.Context) {
	info, queryErr := db.GetLandingPageGeneralInfo(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get landing page general info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, info)
}

func GetHomeGeneralInfo(ctx *gin.Context) {
	info, queryErr := db.GetLandingPageInfo(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get landing page general info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, info)
}

func UpdateHomeGeneralInfo(ctx *gin.Context) {
	data := dto.UpdateLandingPageGeneralInfo_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateLandingPageGeneralInfo(db.DB, db.LandingPageGeneralData{
		Title:                 data.Title,
		StarterSentence:       data.StarterSentence,
		SecondStarterSentence: data.SecondStarterSentence,
		PlansParagraph:        data.PlansParagraph,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to update landing page general info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetPlanParagrarph(ctx *gin.Context) {
	info, queryErr := db.GetPlansParagraph(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get plans paragraph info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, info)
}

func UpdatePlanParagraph(ctx *gin.Context) {
	var text string

	bindErr := ctx.ShouldBindJSON(&text)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdatePlansParagraph(db.DB, text)
	if queryErr != nil {
		common.Logger.Printf("Failed to update plans paragraph info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetHomePlans(ctx *gin.Context) {
	data, queryErr := db.GetPlans(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get all plans info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func GetPlanByID(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "id")
	if params == nil {
		return
	}

	id := params[0]

	plan, queryErr := db.GetPlanByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to get a plan info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, plan)
}

func CreatePlan(ctx *gin.Context) {
	data := dto.CreatePlan_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreatePlan(db.DB, db.Plan{
		Title:       data.Title,
		Description: data.Description,
		Duration:    data.Duration,
		Price:       data.Price,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to create a plan: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func DeletePlanByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "invalid paramter: id")
		return
	}

	queryErr := db.DeletePlanByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a plan by ID (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func ReplacePlanByID(ctx *gin.Context) {
	data := dto.ReplacePlan_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", data)
		return
	}

	queryErr := db.ReplacePlan(db.DB, db.Plan{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Duration:    data.Duration,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to replace a plan: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetAdsInfo(ctx *gin.Context) {
	info, queryErr := db.GetAdsInfo(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed get ads info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, info)
}

func UpdateAdsInfo(ctx *gin.Context) {
	data := dto.UpdateAdsInfo_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateAdsInfo(db.DB, db.AdsInfo{
		Title:       data.Title,
		Description: data.Description,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed update ads info: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetHomeProducts(ctx *gin.Context) {
	products, queryErr := db.GetProducts(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get products: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func GetProductByID(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "id")
	if params == nil {
		return
	}

	id := params[0]

	product, queryErr := db.GetProductByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to get a product by ID (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func CreateHomeProduct(ctx *gin.Context) {
	data := dto.CreateProduct_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreateProduct(db.DB, db.Product{
		Name:        data.Name,
		Description: data.Description,
		Marka:       data.Marka,
		Price:       data.Price,
		CategoryID:  data.CategoryID,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to create a product: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}
func DeleteHomeProduct(ctx *gin.Context) {}
func UpdateHomeProduct(ctx *gin.Context) {
	data := dto.UpdateProduct_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateProduct(db.DB, db.Product{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Marka:       data.Marka,
		Price:       data.Price,
		CategoryID:  data.CategoryID,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to update a product: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func DeleteHomeProductByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: id")
		return
	}

	queryErr := db.DeleteProductByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a product by ID (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetProductCategories(ctx *gin.Context) {
	categories, queryErr := db.GetProductCategories(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get product categories: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func GetCategoryProducts(ctx *gin.Context) {
	catProd, queryErr := db.GetProductCategoriesWithProducts(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get product categories with products: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, catProd)
}

func GetProductsOfCategory(ctx *gin.Context) {
	name := ctx.Query("name")
	idStr := ctx.Query("id")

	if name == "" && idStr == "" {
		ctx.String(http.StatusBadRequest, "One of the two query parameters is required: 'name' or 'id'")
		return
	}

	var category *db.ProductCategory
	var catQueryErr error

	if idStr != "" {
		id, convErr := strconv.ParseInt(idStr, 10, 64)
		if convErr != nil {
			ctx.String(http.StatusBadRequest, "Invalid query paramter: id")
			return
		}

		category, catQueryErr = db.GetProductCategoryByID(db.DB, id)
	} else {
		category, catQueryErr = db.GetProductCategoryByName(db.DB, name)
	}

	if catQueryErr != nil {
		common.Logger.Printf("Failed to get a product category by name: %v\n", catQueryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	products, queryErr := db.GetProductWithCategoryByID(db.DB, category.ID)
	if queryErr != nil {
		common.Logger.Printf("Failed to get products of category: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func CreateProductCategory(ctx *gin.Context) {
	data := dto.CreateProductCategory_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreateProductCategory(db.DB, data.Name)
	if queryErr != nil {
		common.Logger.Printf("Failed to create a category: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func DeleteProductCategoryByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: id")
		return
	}

	queryErr := db.DeleteProductCategoryByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a category: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func DeleteProductCategoryByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: name")
		return
	}

	cat, catQueryErr := db.GetProductCategoryByName(db.DB, name)
	if catQueryErr != nil {
		common.Logger.Printf("Failed to get a category by name: %v\n", catQueryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	queryErr := db.DeleteProductCategoryByID(db.DB, cat.ID)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a category by name: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func MoveProductToCategory(ctx *gin.Context) {
}

func DeleteProductsOfCategory(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: id")
		return
	}

	queryErr := db.DeleteProductsOfCategoryByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete products of category (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func ProductExistsUnderCategory(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "productId", "categoryId")
	if params == nil {
		return
	}

	productId, categoryId := params[0], params[1]

	exists, queryErr := db.ProductExistsUnderCategory(db.DB, productId, categoryId)
	if queryErr != nil {
		common.Logger.Printf("Failed to check if a product exists under a category: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, exists)
}

func GetContacts(ctx *gin.Context) {
	contacts, queryErr := db.GetContacts(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get contacts: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

func UpdateContacts(ctx *gin.Context) {
	contacts := db.Contacts{}

	bindErr := ctx.ShouldBindJSON(&contacts)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateContacts(db.DB, contacts)
	if queryErr != nil {
		common.Logger.Printf("Failed to update contacts: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetQNA(ctx *gin.Context) {
	array, queryErr := db.GetQNA(db.DB)
	if queryErr != nil {
		common.Logger.Printf("failed to get QNA: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, array)
}

func AddQNA(ctx *gin.Context) {
	data := dto.CreateQNA_Req{}
	bindErr := ctx.ShouldBindBodyWithJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.AddQNA(db.DB, data.Question, data.Answer)
	if queryErr != nil {
		common.Logger.Printf("failed to add QNA: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func DeleteQNA(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: id")
		return
	}

	queryErr := db.DeleteQNAByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("failed to delete QNA: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
